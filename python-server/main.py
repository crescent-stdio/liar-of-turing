from fastapi import FastAPI
from openai import OpenAI
from dotenv import load_dotenv
import os
from message_data import MessageData

load_dotenv()
client = OpenAI(api_key = os.getenv("OPENAI_API_KEY"))
app = FastAPI()

@app.post("/useGPT")
def get_message_from_GPT(messageRequest: MessageData):
    print(messageRequest.message)
    chatlog = messageRequest.message

    completion = client.chat.completions.create(
    model="gpt-4",
    # model="gpt-3.5-turbo",
    messages=[
        {"role": "system", "content": '''
        문제: 단체 채팅방의 참여자로서 다음에 이어질 문장은 무엇인가? 모든 조건을 반드시 만족하도록 할 것.
        조건1. 최대 10단어로, 가급적 짧게 말할 것.
        조건2. 주어지는 채팅들과 최대한 비슷한 스타일로 말할 것.
        조건3. 대화의 주제에서 벗어난 단어를 사용하지 말 것.
        조건4. 콤마, 점를 가급적 사용하지 말 것. 맞춤법과 띄어쓰기를 가끔 틀릴 것.
        조건5. 닉네임은 일반적으로 내용과 무관함.
        조건6. 말할 차례를 지킬 것.
        '''},
        {"role": "user", "content": chatlog}])
        
    answer = completion.choices[0].message.content
    if len(answer.split(':')) > 1:
        answer = answer.split(':')[1].lstrip()
    print("result_message: ", answer)
    return {"user_name": messageRequest.user_name, "message": answer}

#     return completion.choices[0].message.content

# print(get_message_from_GPT())
    
