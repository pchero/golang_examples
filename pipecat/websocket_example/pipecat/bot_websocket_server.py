#
# Copyright (c) 2024–2025, Daily
#
# SPDX-License-Identifier: BSD 2-Clause License
#

import os

from loguru import logger
from pipecat.audio.vad.silero import SileroVADAnalyzer
from pipecat.frames.frames import LLMRunFrame
from pipecat.pipeline.pipeline import Pipeline
from pipecat.pipeline.runner import PipelineRunner
from pipecat.pipeline.task import PipelineParams, PipelineTask
from pipecat.processors.aggregators.openai_llm_context import OpenAILLMContext
from pipecat.processors.frameworks.rtvi import RTVIConfig, RTVIObserver, RTVIProcessor
from pipecat.serializers.protobuf import ProtobufFrameSerializer
from pipecat.services.cartesia.tts import CartesiaTTSService
from pipecat.services.deepgram.stt import DeepgramSTTService
from pipecat.services.openai.llm import OpenAILLMService
from pipecat.processors.frame_processor import FrameProcessor

# from pipecat.services.gemini_multimodal_live import GeminiMultimodalLiveLLMService
from pipecat.transports.websocket.server import (
    WebsocketServerParams,
    WebsocketServerTransport,
)

import factory_tts

async def run_bot_websocket_server(args):
    ws_transport = WebsocketServerTransport(
        params=WebsocketServerParams(
            serializer=ProtobufFrameSerializer(),
            audio_in_enabled=True,
            audio_out_enabled=True,
            add_wav_header=False,
            vad_analyzer=SileroVADAnalyzer(),
            session_timeout=60 * 3,  # 3 minutes
        )
    )
    
    stt = factory_tts.create_stt_service(args.stt)

    tts = factory_tts.create_tts_service(
        args.tts,
        voice_id="71a7ad14-091c-4e8e-a314-022ece01c121"
    )

    llm = factory_tts.create_llm_server(args.llm)

    context_aggregator = factory_tts.create_context_and_aggregator(llm, args.messages_file)

    # RTVI events for Pipecat client UI
    rtvi = RTVIProcessor(config=RTVIConfig(config=[]))

    pipeline = Pipeline(
        [
            ws_transport.input(),
            rtvi,  # RTVI processor
            stt,
            context_aggregator.user(),  # User responses
            llm,  # LLM
            tts,  # TTS            ws_transport.output(),
            context_aggregator.assistant(),
            ws_transport.output(),
        ]
    )

    task = PipelineTask(
        pipeline,
        params=PipelineParams(
            enable_metrics=True,
            enable_usage_metrics=True,
        ),
        observers=[RTVIObserver(rtvi)],
    )

    @rtvi.event_handler("on_client_ready")
    async def on_client_ready(rtvi):
        logger.info("Pipecat client ready.")
        await rtvi.set_bot_ready()
        # Kick off the conversation.
        await task.queue_frames([LLMRunFrame()])

    @ws_transport.event_handler("on_client_connected")
    async def on_client_connected(transport, client):
        logger.info("Pipecat Client connected")

    @ws_transport.event_handler("on_client_disconnected")
    async def on_client_disconnected(transport, client):
        logger.info("Pipecat Client disconnected")
        await task.cancel()

    @ws_transport.event_handler("on_session_timeout")
    async def on_session_timeout(transport, client):
        logger.info(f"Entering in timeout for {client.remote_address}")
        await task.cancel()

    runner = PipelineRunner()

    await runner.run(task)
