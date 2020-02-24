ARG ALPINE_VERSION=3.11
FROM golang:1.13-alpine${ALPINE_VERSION} AS builder

COPY . /go/src/github.com/Brian-Williams/ami-share/
WORKDIR /go/src/github.com/Brian-Williams/ami-share/
RUN go build -o ami-share ./

FROM alpine:${ALPINE_VERSION}

COPY LICENSE README.md /

COPY entrypoint.sh /entrypoint.sh
COPY --from=builder ami-share .

ENTRYPOINT ["/entrypoint.sh"]
