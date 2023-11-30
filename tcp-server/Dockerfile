# Builder

FROM golang:1.16-alpine AS builder
RUN apk add --update make git curl gcc

ARG GITHUB_PATH=github.com/pullya/wow_tcp_server/tcp-server

COPY Makefile /home/${GITHUB_PATH}/Makefile
COPY go.mod /home/${GITHUB_PATH}/go.mod
COPY go.sum /home/${GITHUB_PATH}/go.sum
COPY ./ ./

WORKDIR /home/${GITHUB_PATH}
COPY . /home/${GITHUB_PATH}
RUN make build

# TCP Server

FROM alpine:latest as server
LABEL org.opencontainers.image.source https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

ARG GITHUB_PATH=github.com/pullya/wow_tcp_server/tcp-server

# COPY --from=builder /home/${GITHUB_PATH}/bin/tcp-server .

# RUN chown root:root tcp-server

EXPOSE 8081

CMD ["./tcp-server"]