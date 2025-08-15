FROM golang:1.23

WORKDIR /app

COPY . . 

WORKDIR /app/app-gateway

RUN go mod tidy

CMD ["go", "run", "cmd/app-gateway/main.go"]