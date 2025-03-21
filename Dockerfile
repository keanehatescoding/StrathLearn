FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go get github.com/google/uuid

COPY . .

RUN apk add --no-cache gcc musl-dev
RUN go build -o app backend/main.go

EXPOSE 8080

# Make sure to use $PORT environment variable for Heroku
CMD ["sh", "-c", "./app"]