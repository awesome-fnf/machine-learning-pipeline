FROM tensorflow/serving

RUN apt-get update
RUN apt-get -y install python3
RUN apt-get -y install python3-pip
RUN pip3 install oss2
RUN mkdir serving
COPY ./serve.py /serving/serve.py
WORKDIR /serving

ENTRYPOINT ["python3", "serve.py"]