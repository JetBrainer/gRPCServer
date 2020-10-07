FROM golang:1.15.2-alpine

ENV GO111MODULE=on

WORKDIR /opt/code/
ADD ./ /opt/code/


RUN apk update && apk upgrade && \
    apk add --no-cache git

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN cd cmd/grpcserver && go build -o ../../bin/main && cd ../../

ENTRYPOINT ["./bin/main"]