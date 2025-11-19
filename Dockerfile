# -----------------------------
# Stage 1 — Build the Go binary
# -----------------------------
FROM golang:1.24 AS builder

# Create working directory
WORKDIR /app

# Copy source code
COPY . .

# Build statically for Linux amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o worker .

# -----------------------------
# Stage 2 — Minimal runtime image
# -----------------------------
FROM ubuntu:22.04

# Create directory where the app will live
WORKDIR /usr/local/bin

# Copy the built binary from the builder
COPY --from=builder /app/worker .

RUN apt-get update && \
    apt-get install -y build-essential openjdk-11-jdk golang && \
    apt-get clean

# Command to run the binary
ENTRYPOINT ["./worker"]
