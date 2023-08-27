# syntax=docker/dockerfile:1

ARG GO_VERSION=1.21
ARG DISTROLESS=static-debian11:nonroot-amd64

FROM golang:${GO_VERSION} AS base
WORKDIR /usr/src

COPY go.* .
RUN go mod download

COPY . .

FROM base AS test
RUN go test -v -cover -count 1 ./...

FROM base AS build
RUN GOOS=linux GOARCH=amd64 go build \
    -ldflags "-w -s" \
    -o /usr/bin/a ./cmd/spangle-paste/*.go

FROM gcr.io/distroless/${DISTROLESS} AS final
WORKDIR /opt

USER nonroot:nonroot
COPY --from=build --chown=nonroot:nonroot /usr/bin/a .

EXPOSE 8000
CMD ["./a"]