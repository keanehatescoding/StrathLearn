FROM golang:1.22.6-alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/backend
RUN go mod download
RUN go build -o strathlearn .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/backend/strathlearn /app/
COPY ./frontend /app/frontend
COPY ./backend/challenges /app/backend/challenges

RUN apk --no-cache add docker-cli

EXPOSE 8080
CMD ["/app/strathlearn"]
