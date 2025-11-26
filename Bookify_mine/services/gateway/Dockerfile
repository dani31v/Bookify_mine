# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24-alpine

FROM golang:${GO_VERSION} AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /opt/service ./services/gateway/cmd/gateway

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /opt/service /usr/local/bin/service
ENTRYPOINT ["/usr/local/bin/service"]
