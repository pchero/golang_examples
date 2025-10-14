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

import uvicorn
from dotenv import load_dotenv
from fastapi import FastAPI, Request, WebSocket
from fastapi.middleware.cors import CORSMiddleware

# Load environment variables
load_dotenv(override=True)

# from bot_fast_api import run_bot
from bot_websocket_server import run_bot_websocket_server


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Handles FastAPI startup and shutdown."""
    yield  # Run app

parser = argparse.ArgumentParser()
parser.add_argument("--messages_file", type=str, required=True)
parser.add_argument("--llm", type=str, required=True)
parser.add_argument("--tts", type=str, required=True)
parser.add_argument("--stt", type=str, required=True)
args = parser.parse_args()

with open(args.messages_file, "r", encoding="utf-8") as f:
    messages = json.load(f)

print("✅ Pipecat 실행 설정")
print("Messages:", messages)
print("LLM:", args.llm)
print("TTS:", args.tts)
print("STT:", args.stt)


# Initialize FastAPI app with lifespan manager
app = FastAPI(lifespan=lifespan)

# Configure CORS to allow requests from any origin
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.post("/connect")
async def bot_connect(request: Request) -> Dict[Any, Any]:
    ws_url = "ws://localhost:8765"
    return {"ws_url": ws_url}


async def main():
    tasks = []
    try:
        tasks.append(run_bot_websocket_server(args))

        config = uvicorn.Config(app, host="0.0.0.0", port=7860)
        server = uvicorn.Server(config)
        tasks.append(server.serve())

        await asyncio.gather(*tasks)
    except asyncio.CancelledError:
        print("Tasks cancelled (probably due to shutdown).")


if __name__ == "__main__":
    asyncio.run(main())