# Устанавливаем базовое изображение
FROM golang:1.21-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем приложение
RUN go build -o api-gateway .

# Устанавливаем порт
EXPOSE 8080

# Запуск приложения
CMD ["./api-gateway"]
