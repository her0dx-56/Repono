
# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go backend
RUN CGO_ENABLED=0 GOOS=linux go build -o repono ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Copy compiled application
COPY --from=builder /app/repono .

# Backend listens on this port by default
ENV HOST=0.0.0.0
ENV PORT=9000

EXPOSE 9000

CMD ["./repono"]
