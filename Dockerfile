FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go build -o weather-service main.go

CMD ["./weather-service"]