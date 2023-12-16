from pydantic import BaseModel

class MessageData(BaseModel):
    gptName: str
    message: str