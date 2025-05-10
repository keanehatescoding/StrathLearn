FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o strathlearn ./backend

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/strathlearn /app/
COPY ./frontend /app/frontend
COPY ./backend/challenges /app/backend/challenges

RUN apk --no-cache add docker-cli

EXPOSE 8080
CMD ["/app/strathlearn"]