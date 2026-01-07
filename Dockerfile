# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o pipeops .

# Final stage
FROM scratch

# Copy certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/pipeops /usr/local/bin/pipeops

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/pipeops"]

# Default command
CMD ["--help"]

# Labels
LABEL org.opencontainers.image.title="PipeOps CLI"
LABEL org.opencontainers.image.description="Official command line tool for PipeOps"
LABEL org.opencontainers.image.url="https://github.com/PipeOpsHQ/pipeops-cli"
LABEL org.opencontainers.image.documentation="https://docs.pipeops.io"
LABEL org.opencontainers.image.source="https://github.com/PipeOpsHQ/pipeops-cli"
LABEL org.opencontainers.image.vendor="PipeOps"
LABEL org.opencontainers.image.licenses="MIT"