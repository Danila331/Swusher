# filepath: c:\Users\danil\OneDrive\Рабочий стол\ShareHub\back-end\Dockerfile
FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod tidy

CMD ["go", "run", "cmd/app/main.go"]