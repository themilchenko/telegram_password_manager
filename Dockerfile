# syntax=docker/dockerfile:1

FROM ubuntu:latest
WORKDIR /app
COPY . .

RUN apt-get update && apt-get install -y golang libssl-dev ca-certificates
RUN update-ca-certificates
RUN go mod download
RUN go build -o main ./cmd/main.go

CMD ["./main"]
