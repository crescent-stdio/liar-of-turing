from pydantic import BaseModel

class MessageData(BaseModel):
    user_UUID: str
    message: str