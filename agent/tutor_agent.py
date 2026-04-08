"""
FluentFlow LiveKit voice tutor (MVP).
Dispatched when the learner connects (RoomAgentDispatch embedded in join token).

Requires OPENAI_API_KEY for the realtime speech pipeline. Without it, the process
stays up so `docker compose` still succeeds — add the key for real voice.
"""
from __future__ import annotations

import asyncio
import json
import logging
import os
import sys

from dotenv import load_dotenv

from livekit import agents
from livekit.agents import Agent, AgentServer, AgentSession
from livekit.plugins import openai, silero

load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("fluentflow-tutor")

AGENT_NAME = os.getenv("LIVEKIT_AGENT_NAME", "fluentflow-tutor")
REALTIME_MODEL = os.getenv("OPENAI_REALTIME_MODEL", "gpt-4o-realtime-preview")
TRANSCRIPTION_MODEL = os.getenv("OPENAI_TRANSCRIPTION_MODEL", "gpt-4o-mini-transcribe")
TTS_VOICE = os.getenv("OPENAI_TTS_VOICE", "alloy")


def scenario_prompt(scenario_id: str, meta: dict) -> str:
    target = meta.get("target_language", "the target language")
    level = meta.get("proficiency_level", "intermediate")
    goal = meta.get("learning_goal", "general practice")
    return f"""You are FluentFlow, a supportive language tutor for scenario "{scenario_id}".
The learner's stated level is {level}. Learning goal: {goal}.
They are practicing speaking {target}.
Speak mostly in {target}; briefly use English only when explaining a correction.
Keep replies concise (one or two short sentences) so the conversation stays natural.
Correct major errors gently after responding to what they said.
Do not use markdown, emojis, or bullet lists in speech."""


class Tutor(Agent):
    def __init__(self, instructions: str) -> None:
        super().__init__(instructions=instructions)


server = AgentServer()


def build_realtime_model() -> openai.realtime.RealtimeModel:
    """Build RealtimeModel with explicit GPT model + STT config, fallback safely if SDK args differ."""
    try:
        return openai.realtime.RealtimeModel(
            model=REALTIME_MODEL,
            voice=TTS_VOICE,
            input_audio_transcription={"model": TRANSCRIPTION_MODEL},
        )
    except TypeError:
        logger.warning(
            "RealtimeModel kwargs unsupported by installed SDK; using fallback constructor",
            exc_info=True,
        )
        return openai.realtime.RealtimeModel(voice=TTS_VOICE)


@server.rtc_session(agent_name=AGENT_NAME)
async def entrypoint(ctx: agents.JobContext) -> None:
    meta: dict = {}
    raw_meta = ""
    if getattr(ctx, "job", None) is not None and getattr(ctx.job, "metadata", ""):
        raw_meta = ctx.job.metadata
    if raw_meta:
        try:
            meta = json.loads(raw_meta)
        except json.JSONDecodeError:
            logger.warning("invalid job metadata JSON")

    scenario_id = meta.get("scenario_id", "practice")
    instructions = scenario_prompt(scenario_id, meta)

    if not os.getenv("OPENAI_API_KEY"):
        logger.warning(
            "OPENAI_API_KEY is not set — voice pipeline disabled. "
            "Set it in .env for OpenAI Realtime."
        )
        await ctx.connect()
        return

    logger.info(
        "starting realtime tutor model=%s stt=%s voice=%s",
        REALTIME_MODEL,
        TRANSCRIPTION_MODEL,
        TTS_VOICE,
    )
    session = AgentSession(vad=silero.VAD.load(), llm=build_realtime_model())

    await session.start(room=ctx.room, agent=Tutor(instructions))

    greeting = (
        f"Greet the learner briefly in the practice scenario '{scenario_id}' "
        "and ask one simple opening question."
    )
    # Retry the initial turn: Realtime can cancel if it thinks user audio started.
    for attempt in range(1, 4):
        try:
            await session.generate_reply(instructions=greeting)
            logger.info("initial greeting generated on attempt %d", attempt)
            break
        except TypeError:
            # Older SDKs may not support extra kwargs; retry plain call.
            await session.generate_reply()
            logger.info("initial greeting generated with plain call")
            break
        except Exception:
            logger.warning("initial greeting attempt %d failed", attempt, exc_info=True)
            await asyncio.sleep(0.8)


if __name__ == "__main__":
    agents.cli.run_app(server)
