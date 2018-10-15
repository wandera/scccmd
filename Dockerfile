# Builder image
FROM golang:1.11 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v

# Runtime image
FROM alpine:3.8
RUN apk --no-cache add ca-certificates

COPY --from=builder /build/scccmd /app/scccmd
WORKDIR /app

ENTRYPOINT ["./scccmd"]