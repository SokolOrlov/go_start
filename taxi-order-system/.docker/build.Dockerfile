FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY internal ./internal
COPY pkg ./pkg
COPY cmd ./cmd
COPY .config ./config
COPY migrations ./migrations

WORKDIR /app/cmd/client
RUN go build -o app

WORKDIR /app/cmd/trip
RUN go build -o app

WORKDIR /app/cmd/driver
RUN go build -o app