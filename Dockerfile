FROM golang:1.24.2 AS builder
LABEL maintainer="Patrick Hermann patrick.hermann@sva.de"
LABEL org.opencontainers.image.source https://github.com/stuttgart-things/kaeffken

WORKDIR app

COPY . .

RUN go mod tidy
RUN go build -o /bin/kaeffken

FROM alpine:3.21.3

COPY --from=builder /bin/kaeffken /usr/bin/kaeffken

ENTRYPOINT ["kaeffken"]
