FROM golang:alpine AS build-env

ADD . /go/src/rsyslog_exporter
WORKDIR /go/src/rsyslog_exporter

RUN apk add --no-cache git
RUN go install -v -ldflags "-X main.programVersion=$(git describe --tags || git rev-parse --short HEAD || echo dev)" ./...

FROM alpine:latest
COPY --from=build-env /go/bin/rsyslog_exporter /rsyslog_exporter

ENTRYPOINT ["/rsyslog_exporter"]

EXPOSE 9120
