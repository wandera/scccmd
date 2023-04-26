# Builder image
FROM golang:1.20.3-alpine3.17 AS builder

WORKDIR /build

ARG VERSION

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build GOMODCACHE=/go/pkg/mod GOCACHE=/root/.cache/go-build go build -v -ldflags '-w -s -X 'github.com/wandera/scccmd/cmd.Version=${VERSION}

# Docker Cloud args, from hooks/build.
ARG CACHE_TAG
ENV CACHE_TAG ${CACHE_TAG}

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags '-w -s -X 'github.com/wandera/scccmd/cmd.Version=${CACHE_TAG}

# Runtime image
FROM alpine:3.17.3
RUN apk --no-cache add ca-certificates

COPY --from=builder /build/scccmd /app/scccmd
WORKDIR /app

ENTRYPOINT ["./scccmd"]
