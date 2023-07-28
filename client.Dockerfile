FROM python:3.11-slim-bookworm

RUN apt update && apt install -y build-essential libssl-dev

WORKDIR /app

COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

COPY pow.py pow.py
COPY client.py client.py

CMD ["python3", "client.py"]
