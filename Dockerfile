FROM golang:1.21-alpine as builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy Go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main .

# Production stage
FROM alpine:latest

# Install Docker client
RUN apk add --no-cache docker-cli

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/main .

# Create necessary directories
RUN mkdir -p ./challenges ./executor-workdir ./frontend

# Copy frontend and challenge files if they exist
COPY --from=builder /build/frontend ./frontend
COPY --from=builder /build/challenges ./challenges

EXPOSE 8080 9000

CMD ["./main"]