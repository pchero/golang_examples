import asyncio
import os
import json

from loguru import logger

# tts
from pipecat.services.cartesia.tts import CartesiaTTSService
from pipecat.services.elevenlabs.tts import ElevenLabsTTSService
from pipecat.services.google.tts import GoogleTTSService

# stt
from pipecat.services.deepgram.stt import DeepgramSTTService
from pipecat.services.whisper.stt import Model, WhisperSTTService
from pipecat.processors.filters.stt_mute_filter import STTMuteConfig, STTMuteFilter, STTMuteStrategy

# llm
from pipecat.services.openai.llm import OpenAILLMService

# aggregators / context
from pipecat.processors.aggregators.openai_llm_context import OpenAILLMContext

# pipeline
from pipecat.audio.vad.silero import SileroVADAnalyzer
from pipecat.frames.frames import LLMRunFrame
from pipecat.pipeline.pipeline import Pipeline
from pipecat.pipeline.runner import PipelineRunner
from pipecat.pipeline.task import PipelineParams, PipelineTask
from pipecat.processors.aggregators.openai_llm_context import OpenAILLMContext
from pipecat.processors.frameworks.rtvi import RTVIConfig, RTVIObserver, RTVIProcessor
from pipecat.serializers.protobuf import ProtobufFrameSerializer
from pipecat.transports.websocket.client import (
    WebsocketClientParams,
    WebsocketClientTransport,
)

async def run_pipeline(args, initial_messages):
    logger.info(f"Connecting Pipecat client to Go WebSocket server at: {args.ws_server_url}")

    ws_transport = WebsocketClientTransport(
        uri=args.ws_server_url,
        params=WebsocketClientParams(
            serializer=ProtobufFrameSerializer(),
            audio_in_enabled=True,
            audio_out_enabled=True,
            add_wav_header=False,
            vad_analyzer=SileroVADAnalyzer(),
            session_timeout=60 * 3,
        )
    )

    stt = create_stt_service(args.stt)

    tts = create_tts_service(
        args.tts,
        voice_id="71a7ad14-091c-4e8e-a314-022ece01c121"
    )

    llm = create_llm_server(args.llm)

    context_aggregator = llm.create_context_aggregator(OpenAILLMContext(
        messages=initial_messages
    ))

    rtvi = RTVIProcessor(config=RTVIConfig(config=[]))

    pipeline = Pipeline(
        [
            ws_transport.input(),
            rtvi,
            stt,
            context_aggregator.user(),
            llm,
            tts,
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


def create_tts_service(name: str, **options):
    """
    Factory function to create a Pipecat TTS service instance
    based on the service name and optional parameters.
    """
    name = name.lower()

    if name == "cartesia":
        return CartesiaTTSService(
            api_key=options.get("api_key", os.getenv("CARTESIA_API_KEY")),
            voice_id=options.get("voice_id", "default_voice_id")
        )
    elif name == "elevenlabs":
        return ElevenLabsTTSService(
            api_key=options.get("api_key", os.getenv("ELEVENLABS_API_KEY")),
            voice_id=options.get("voice_id", "default_voice_id")
        )
    else:
        raise ValueError(f"Unsupported TTS service: {name}")


def create_stt_service(name: str, **options):
    """
    Factory function to create a Pipecat STT service instance
    based on the service name and optional parameters.
    """
    name = name.lower()

    if name == "deepgram":
        return DeepgramSTTService(
            api_key=options.get("api_key", os.getenv("DEEPGRAM_API_KEY"))
        )
    else:
        raise ValueError(f"Unsupported STT service: {name}")


def create_llm_server(name: str, **options):
    """
    Factory function to create a Pipecat LLM service instance
    based on the argument in 'service.model' format.
    """
    if "." not in name:
        raise ValueError(f"Wrong LLM: {name}. LLM argument must be in 'service.model' format, e.g., 'openai.gpt-4o-mini'")

    service_name, model_name = name.split(".", 1)
    service_name = service_name.lower()

    if service_name == "openai":
        llm = OpenAILLMService(
            api_key=options.get("api_key", os.getenv("OPENAI_API_KEY")),
            model=model_name
        )
        llm.context_class = OpenAILLMContext
    else:
        raise ValueError(f"Unsupported LLM service: {service_name}")
    
    return llm


def create_context_and_aggregator(llm, file_path: str):
    print("Executing create_context_and_aggregator. file_path:", file_path)
    
    with open(file_path, "r", encoding="utf-8") as f:
        messages = json.load(f)

    for msg in messages:
        if "role" not in msg or "content" not in msg:
            raise ValueError("Each message must contain 'role' and 'content'")


    context = llm.context_class(messages)
    context_aggregator = llm.create_context_aggregator(context)
    
    return context_aggregator