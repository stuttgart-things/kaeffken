FROM golang:1.23.4 AS builder
LABEL maintainer="Patrick Hermann patrick.hermann@sva.de"
LABEL org.opencontainers.image.source https://github.com/stuttgart-things/kaeffken

WORKDIR app

COPY . .

RUN go mod tidy
RUN go build -o /bin/kaeffken

FROM alpine:3.19.1

COPY --from=builder /bin/kaeffken /usr/bin/kaeffken

ENTRYPOINT ["kaeffken"]
