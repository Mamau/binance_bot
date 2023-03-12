FROM golang:1.19-alpine

ARG CGO_ENABLED=0

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download && go install github.com/cosmtrek/air@latest
