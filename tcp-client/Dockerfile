# Builder

FROM golang:1.21-alpine AS builder
WORKDIR /usr/local/src

RUN apk --no-cache add bash make git curl gcc musl-dev

COPY Makefile ./
COPY go.mod ./

CMD go mod download

COPY . ./
RUN go build -o ./bin/tcp-client cmd/tcp-client/main.go

# TCP-client

FROM alpine:latest AS client
COPY --from=builder /usr/local/src/bin/tcp-client /

CMD ["/tcp-client"]