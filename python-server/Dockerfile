# 베이스 이미지 선택
FROM python:3.9-alpine

# 작업 디렉토리 설정
WORKDIR /app

# 의존성 파일 복사 및 설치
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# 소스 코드 복사
COPY . .

# 컨테이너 실행 시 실행할 명령어
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "5000"]
