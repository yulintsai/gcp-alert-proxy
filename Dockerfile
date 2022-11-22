# build stage
FROM golang:1.18.8-alpine3.16 AS builder

LABEL stage=gcp-alert-proxy-intermediate

ENV GO111MODULE=on

ADD ./ /go/src/gcp-alert-proxy

RUN cd /go/src/gcp-alert-proxy && go build -mod vendor

FROM alpine:3.16.3

RUN apk add --no-cache tzdata

COPY --from=builder /go/src/gcp-alert-proxy/gcp-alert-proxy ./

