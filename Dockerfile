FROM golang:1.26.3 AS builder

RUN update-ca-certificates

WORKDIR /app

RUN apt-get update -y && apt-get install -y libxml2-dev pkg-config
RUN go install golang.org/x/lint/golint@latest
COPY . .
RUN go build -o eadindexer

ENTRYPOINT [ "./eadindexer" ]
