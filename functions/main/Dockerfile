FROM golang:1.12.7-stretch

ENV GO111MODULE on
RUN go get k8s.io/client-go@kubernetes-1.15.0
RUN go get github.com/aliyun/alibaba-cloud-sdk-go/services/cs
RUN go get github.com/awesome-fc/golang-runtime
RUN mkdir -p /go/src/fnf-fc-golang-k8s-builder
COPY *.go /go/src/fnf-fc-golang-k8s-builder/
COPY *.mod /go/src/fnf-fc-golang-k8s-builder/
RUN cd /go/src/fnf-fc-golang-k8s-builder; go build -o bin/bootstrap /go/src/fnf-fc-golang-k8s-builder/*.go
