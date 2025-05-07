# STEP 1: Build
FROM golang:1.22 AS builder

WORKDIR /app

# Copy the entire project
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bitcoin-tracker ./cmd/main.go

# STEP 2: Runtime
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bitcoin-tracker .
COPY .env .

# Add user and install curl for healthchecks
RUN adduser -S -D -H -h /app btcuser && \
    chown -R btcuser: /app && \
    apk add --no-cache curl

# Switch to non-root user
USER btcuser

# Expose web port
EXPOSE 8080

# Run the app
CMD ["./bitcoin-tracker"]
