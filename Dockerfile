FROM golang:1.14 as builder
WORKDIR /go/src/github.com/rolandvarga/batch_proc/
COPY . .
RUN make all

FROM alpine:latest
COPY --from=builder /go/src/github.com/rolandvarga/batch_proc/batch_proc /go/src/github.com/rolandvarga/batch_proc/
COPY --from=builder /go/src/github.com/rolandvarga/batch_proc/data.json /go/src/github.com/rolandvarga/batch_proc/

ENV GOPATH /go/
WORKDIR /go/src/github.com/rolandvarga/batch_proc/

LABEL maintainer="Roland Varga <roland.varga.can@gmail.com>" \
    "org.label-schema.name"="batch_proc" \
    "org.label-schema.base-image.name"="docker.io/library/alpine" \
    "org.label-schema.base-image.version"="latest" \
    "org.label-schema.description"="batch_proc in a container" \
    "org.label-schema.url"="https://github.com/rolandvarga/batch_proc" \
    "org.label-schema.vcs-url"="https://github.com/rolandvarga/batch_proc" \
    "org.label-schema.vendor"="rolandvarga" \
    "org.label-schema.schema-version"="1.0.0" \
    "org.label-schema.usage"="Please see README.md"