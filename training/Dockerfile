FROM docker.io/tensorflow/tensorflow:2.1.0rc1-py3

RUN pip install matplotlib
RUN pip install requests
RUN pip install oss2
RUN mkdir training
COPY ./train.py /training/train.py
WORKDIR /training