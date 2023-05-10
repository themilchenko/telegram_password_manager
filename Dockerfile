# syntax=docker/dockerfile:1

FROM golang:latest
WORKDIR /app
COPY . .

RUN apt-get update && apt-get install -y libssl1.1-dev

RUN go mod tidy
RUN go build -o main ./cmd/main.go

CMD ["./main"]