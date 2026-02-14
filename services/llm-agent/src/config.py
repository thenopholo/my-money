import os
from pathlib import Path

from dotenv import load_dotenv

# Carrega .env do root do projeto (2 níveis acima de services/llm-agent/)
_project_root = Path(__file__).resolve().parent.parent.parent.parent
_env_file = _project_root / ".env"
load_dotenv(_env_file)

OPENAI_API_KEY = os.getenv("OPENAI_API_KEY", "")
LLM_AGENT_PORT = int(os.getenv("LLM_AGENT_PORT", "8001"))
LLM_MODEL = os.getenv("LLM_MODEL", "gpt-5-mini")
