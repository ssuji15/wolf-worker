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
FROM debian:bookworm-slim

# Create directory where the app will live
WORKDIR /usr/local/bin

# Copy the built binary from the builder
COPY --from=builder /app/worker .

RUN apt-get update \
    && apt-get install -y clang \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY ./stdc++.h /usr/include/c++/12/bits/stdc++.h

RUN clang++ -O1 -x c++-header /usr/include/c++/12/bits/stdc++.h -o /usr/include/c++/12/bits/stdc++.h.pch

# Command to run the binary
ENTRYPOINT ["/usr/local/bin/worker"]
