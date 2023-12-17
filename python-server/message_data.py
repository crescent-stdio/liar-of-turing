from pydantic import BaseModel

class MessageData(BaseModel):
    user_name: str
    message: str