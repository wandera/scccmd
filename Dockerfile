# Builder image
FROM golang:1.12.1 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

# Docker Cloud args, from hooks/build.
ARG CACHE_TAG
ENV CACHE_TAG ${CACHE_TAG}

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags '-w -s -X 'github.com/WanderaOrg/scccmd/cmd.Version=${CACHE_TAG}

# Runtime image
FROM alpine:3.8
RUN apk --no-cache add ca-certificates

COPY --from=builder /build/scccmd /app/scccmd
WORKDIR /app

ENTRYPOINT ["./scccmd"]