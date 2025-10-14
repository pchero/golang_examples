# bot_client_pipeline.py (이전 bot_websocket_server.py의 새 이름)
#
# Copyright (c) 2024–2025, Daily
#
# SPDX-License-Identifier: BSD 2-Clause License
#

import asyncio
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

# WebsocketServerTransport 대신 WebsocketClientTransport를 사용합니다.
from pipecat.transports.websocket.client import ( # 클라이언트 트랜스포트 임포트
    WebsocketClientParams,
    WebsocketClientTransport,
)

import factory_tts

# run_bot_websocket_server 대신 run_bot_client_pipeline으로 함수 이름 변경
# Go 서버의 URL과 메시지들을 인자로 받도록 변경
async def run_bot_client_pipeline(args, initial_messages):
    logger.info(f"Connecting Pipecat client to Go WebSocket server at: {args.ws_server_url}")

    # WebsocketServerTransport 대신 WebsocketClientTransport 초기화
    ws_transport = WebsocketClientTransport(
        uri=args.ws_server_url, # Go 서버의 URL을 여기에 전달
        params=WebsocketClientParams(
            serializer=ProtobufFrameSerializer(),
            audio_in_enabled=True,
            audio_out_enabled=True,
            add_wav_header=False,
            vad_analyzer=SileroVADAnalyzer(),
            session_timeout=60 * 3,  # 3 minutes (클라이언트에서는 필요 없을 수도 있지만 유지)
        )
    )

    stt = factory_tts.create_stt_service(args.stt)

    tts = factory_tts.create_tts_service(
        args.tts,
        voice_id="71a7ad14-091c-4e8e-a314-022ece01c121"
    )

    llm = factory_tts.create_llm_server(args.llm)

    # context_aggregator는 이제 initial_messages를 직접 사용할 수 있습니다.
    context_aggregator = llm.create_context_aggregator(OpenAILLMContext(
        # llm=llm,
        messages=initial_messages # 파이썬 클라이언트가 로드한 메시지 전달
    ))
    # context_aggregator = llm.create_context_aggregator(context)


    # RTVI events for Pipecat client UI
    rtvi = RTVIProcessor(config=RTVIConfig(config=[]))

    pipeline = Pipeline(
        [
            ws_transport.input(), # Go 서버로부터 입력
            rtvi,  # RTVI processor
            stt,
            context_aggregator.user(),  # User responses
            llm,  # LLM
            tts,  # TTS
            context_aggregator.assistant(),
            ws_transport.output(), # Go 서버로 출력
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

    # # Go 서버로부터 'client_ready' 또는 유사한 신호를 받을 때까지 기다릴 수 있습니다.
    # # 아니면 파이썬 클라이언트가 연결되면 바로 대화를 시작할 수도 있습니다.
    # @rtvi.event_handler("on_client_ready")
    # async def on_client_ready(rtvi):
    #     logger.info("Pipecat client ready. (Connected to Go server)")
    #     await rtvi.set_bot_ready()
    #     # Go 서버로 봇이 준비되었음을 알릴 수 있습니다.
    #     # await ws_transport.send(json.dumps({"type": "bot_ready", "payload": "ready"}))
    #     # 대화 시작
    #     await task.queue_frames([LLMRunFrame()])


    # on_client_connected/disconnected 이벤트는 WebsocketClientTransport에 맞게 변경되어야 합니다.
    # WebsocketClientTransport는 연결 성공 시 자동으로 처리됩니다.
    # on_client_connected 이벤트는 WebsocketClientTransport에서는 일반적으로 필요 없습니다.
    # on_client_disconnected는 연결이 끊어졌을 때 처리할 수 있습니다.
    @ws_transport.event_handler("on_disconnected") # 클라이언트 트랜스포트의 disconnect 이벤트
    async def on_client_disconnected(transport, error):
        logger.info(f"Pipecat Client disconnected from Go server. Error: {error}")
        await task.cancel()

    @ws_transport.event_handler("on_error") # 클라이언트 트랜스포트의 에러 이벤트
    async def on_error(transport, error):
        logger.error(f"Pipecat Client WebSocket error: {error}")
        await task.cancel()


    # on_session_timeout은 클라이언트 트랜스포트에서는 직접적으로 사용되지 않을 수 있습니다.
    # Go 서버가 타임아웃을 관리하고 연결을 끊을 수 있습니다.

    runner = PipelineRunner()

    # 파이프라인이 준비되고 러너가 시작된 후, 바로 LLM에게 첫 응답을 요청합니다.
    # 이렇게 하면 Go 서버에 연결된 후 바로 웰컴 메시지가 생성됩니다.
    await task.queue_frames([LLMRunFrame()])

    try:
        await runner.run(task)
    except asyncio.CancelledError:
        logger.info("Pipecat client pipeline cancelled.")
    except Exception as e:
        logger.error(f"Pipecat client pipeline error: {e}")