### 1. Update the Dockerfile to reset the entrypoint
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o worker cmd/main.go

FROM linuxserver/ffmpeg:version-7.1-cli

WORKDIR /app
# Copy the binary from builder
COPY --from=builder /app/worker .

# Environment variables for tuning
ENV WORKER_POOL_SIZE=5
ENV USE_GPU=false

# --- CHANGE HERE ---
# Reset the entrypoint so it doesn't run 'ffmpeg' by default
ENTRYPOINT ["/app/worker"]