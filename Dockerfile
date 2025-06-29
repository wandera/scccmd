# Builder image
FROM golang:1.24-alpine3.21 AS builder

WORKDIR /build

ARG VERSION

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build GOMODCACHE=/go/pkg/mod GOCACHE=/root/.cache/go-build go build -v -ldflags '-w -s -X 'github.com/wandera/scccmd/cmd.Version=${VERSION}

# Runtime image
FROM alpine:3.21
RUN apk --no-cache add ca-certificates

COPY --from=builder /build/scccmd /app/scccmd
WORKDIR /app

ENTRYPOINT ["./scccmd"]
