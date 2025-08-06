# filepath: c:\Users\danil\OneDrive\Рабочий стол\ShareHub\back-end\Dockerfile
FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod tidy

CMD ["go", "run", "cmd/app-gateway/main.go"]