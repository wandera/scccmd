# Builder image
FROM golang:1.10 AS builder
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 \
    && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/wanderaorg/scccmd/config
WORKDIR /go/src/github.com/wanderaorg/scccmd/config

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY config/ .
RUN go test -v ./... && CGO_ENABLED=0 go build -o ../bin/config


# Runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/wanderaorg/scccmd/bin/config /app/config
WORKDIR /app

ENTRYPOINT ["./config"]