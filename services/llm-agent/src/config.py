import os

from dotenv import load_dotenv

load_dotenv()

OPENAI_API_KEY = os.getenv("OPENAI_API_KEY", "")
LLM_AGENT_PORT = int(os.getenv("LLM_AGENT_PORT", "8001"))
LLM_MODEL = os.getenv("LLM_MODEL", "gpt-5-mini")
