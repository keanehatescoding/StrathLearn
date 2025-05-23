FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/backend
RUN go mod download
RUN go build -o strathlearn .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/backend/strathlearn /app/

# --- Add this line to copy your data directory ---
COPY ./backend/data /app/data/
# -------------------------------------------------

COPY ./frontend /app/frontend
COPY ./backend/challenges /app/backend/challenges

# This installs Docker CLI in your app image.
# Ensure this is truly needed, as it increases image size and attack surface.
# It's unusual unless your Go app directly interacts with Docker.
RUN apk --no-cache add docker-cli

EXPOSE 8080
CMD ["/app/strathlearn"]