# syntax=docker/dockerfile:1

ARG GO_VERSION=1.21

FROM golang:${GO_VERSION} AS base
WORKDIR /usr/src

COPY go.* .
RUN go mod download

COPY . .

FROM base AS test
RUN go test -v -cover -count 1 ./...

FROM base AS build

ARG LOCATION=todo
RUN GOOS=linux GOARCH=amd64 go build \
    -ldflags "-w -s" \
    -o /usr/bin/a ./cmd/${LOCATION}/*.go