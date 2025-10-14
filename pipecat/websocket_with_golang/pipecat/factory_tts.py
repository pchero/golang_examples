import os
import json

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