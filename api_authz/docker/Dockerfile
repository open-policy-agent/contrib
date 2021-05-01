FROM python:3.9.4-slim-buster

RUN pip install flask requests

COPY echo_server.py /

ENV FLASK_APP=./echo_server.py

CMD ["flask", "run", "--host=0.0.0.0"]

EXPOSE 5000
