# Builder

FROM golang:1.21-alpine AS builder
WORKDIR /usr/local/src

RUN apk --no-cache add bash make git curl gcc musl-dev

COPY Makefile ./
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./
RUN go build -o ./bin/tcp-server cmd/tcp-server/main.go

# TCP Server

FROM alpine:latest AS server
COPY --from=builder /usr/local/src/bin/tcp-server /

CMD ["/tcp-server"]

EXPOSE 8081