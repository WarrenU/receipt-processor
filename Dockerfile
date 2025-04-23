# builder stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Alpine needs musl, but we can still do CGO_ENABLED=0
ENV CGO_ENABLED=0 \
    GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o receipt-processor main.go

# final stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/receipt-processor /receipt-processor

EXPOSE 8080
ENTRYPOINT ["/receipt-processor"]
