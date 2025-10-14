# python_client_for_go.py (이전 server.py의 새 이름)
#
# Copyright (c) 2025, Daily
#
# SPDX-License-Identifier: BSD 2-Clause License
#
import asyncio
import os
from contextlib import asynccontextmanager
from typing import Any, Dict

import argparse
import json
import websockets # 웹소켓 클라이언트 라이브러리 추가

from dotenv import load_dotenv
import sys
from loguru import logger

# Load environment variables
load_dotenv(override=True)

# 기존 run_bot_websocket_server를 run_bot_client_pipeline으로 변경
# 이 함수는 이제 Go 서버에 연결된 WebsocketClientTransport를 사용할 것입니다.
from bot_client_pipeline import run_bot_client_pipeline


async def python_client_main():
    logger.info("--- Received raw arguments ---")
    logger.info(f"sys.argv: {sys.argv}")
    logger.info("----------------------------")

    parser = argparse.ArgumentParser()
    # Go 서버의 웹소켓 주소를 인자로 받도록 변경
    parser.add_argument("--ws_server_url", type=str, required=True,
                        help="WebSocket URL of the Go server (e.g., ws://localhost:8080/ws)")
    parser.add_argument("--llm", type=str, required=True)
    parser.add_argument("--tts", type=str, required=True)
    parser.add_argument("--stt", type=str, required=True)
    # messages_file 인자 추가
    parser.add_argument("--messages_file", type=str, required=True,
                        help="Path to the JSON file containing initial messages for the LLM context.")
    args = parser.parse_args()

    # messages_file에서 메시지 로드
    messages = []
    try:
        with open(args.messages_file, "r", encoding="utf-8") as f:
            messages = json.load(f)
        logger.info(f"✅ Loaded {len(messages)} initial messages from {args.messages_file}")
    except FileNotFoundError:
        logger.info(f"⚠️ Warning: messages_file not found at {args.messages_file}. Starting with empty messages.")
    except json.JSONDecodeError:
        logger.info(f"❌ Error decoding JSON from {args.messages_file}. Starting with empty messages.")


    logger.info("✅ Pipecat Client 설정")
    logger.info(f"Go WebSocket Server URL: {args.ws_server_url}")
    logger.info(f"LLM: {args.llm}")
    logger.info(f"TTS: {args.tts}")
    logger.info(f"STT: {args.stt}")
    # 로드된 메시지 일부 출력 (옵션)
    if messages:
        logger.info(f"Initial Messages (first 2): {messages[:2]}")
    else:
        logger.info("No initial messages loaded.")


    # 이제 run_bot_client_pipeline 함수가 Go 서버의 URL과 messages를 인자로 받아
    # Pipecat WebsocketClientTransport를 초기화하도록 합니다.
    await run_bot_client_pipeline(args, messages)


if __name__ == "__main__":
    try:
        asyncio.run(python_client_main())
    except asyncio.CancelledError:
        logger.info("Python client tasks cancelled.")
    except Exception as e:
        logger.info(f"Python client encountered an error: {e}")