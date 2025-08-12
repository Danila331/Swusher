# Swusher

Swusher — это современная платформа для аренды и обмена вещами между пользователями.  
Проект построен на микросервисной архитектуре с использованием Go, Docker, Prometheus, Grafana и других современных технологий.

---

## 📦 Архитектура проекта

```
Swusher/
├── api/                # Swagger/OpenAPI спецификации
├── app-gateway/        # Основной backend gateway (Go, Echo)
│   ├── cmd/app-gateway # Точка входа приложения
│   └── internal/       # Внутренние пакеты (handlers, middleware, models, servers, static, templates)
│   └── pkg/            # Вспомогательные библиотеки
├── docs/               # Документация
├── infra/
│   └── prometheus/     # Конфигурация Prometheus
├── proto/              # gRPC/Protobuf спецификации для микросервисов
├── services/           # Каталог микросервисов (advertisements, analytics, bookings, chat, media, notification, payments, users)
├── .env                # Переменные окружения
├── Dockerfile          # Dockerfile для сборки app-gateway
├── docker-compose.yml  # Docker Compose для локального запуска
└── README.md           # Описание проекта
```

---

## 🚀 Быстрый старт

1. **Клонируйте репозиторий:**
   ```sh
   git clone https://github.com/Danila331/Swusher.git
   cd Swusher
   ```

2. **Запустите все сервисы через Docker Compose:**
   ```sh
   docker-compose up --build
   ```

3. **Откройте в браузере:**
   - Приложение: [http://localhost:8080](http://localhost:8080)
   - Prometheus: [http://localhost:9090](http://localhost:9090)
   - Grafana: [http://localhost:3000](http://localhost:3000) (логин/пароль: admin/admin)

---

## 🛠️ Технологии

- **Go** — основной язык backend-сервисов
- **Echo** — быстрый и минималистичный web-фреймворк для Go
- **PostgreSQL** — основная база данных
- **Prometheus** — мониторинг и сбор метрик
- **Grafana** — визуализация метрик
- **Docker & Docker Compose** — контейнеризация и оркестрация
- **Swagger/OpenAPI** — документация REST API
- **gRPC/Protobuf** — взаимодействие между микросервисами

---

## 📑 Документация

- **Swagger UI:**  
  Описание REST API находится в [`api/swagger.yaml`](api/swagger.yaml).
- **Prometheus:**  
  Конфигурация в [`infra/prometheus/prometheus.yml`](infra/prometheus/prometheus.yml).

---

## 📂 Основные директории

- `app-gateway/internal/handlers` — обработчики HTTP-запросов
- `app-gateway/internal/midlewary` — middleware (аутентификация, метрики и др.)
- `app-gateway/internal/models` — модели данных
- `app-gateway/internal/servers` — запуск и настройка серверов
- `app-gateway/internal/static` — статические файлы (JS, CSS, изображения)
- `app-gateway/internal/templates` — HTML-шаблоны

---

## 📊 Мониторинг

- Метрики экспонируются на `/metrics` (порт 2112).
- Prometheus собирает метрики, Grafana отображает дашборды.

---

## 🤝 Контрибьютинг

PR и предложения приветствуются!  
Перед отправкой изменений, пожалуйста, убедитесь, что проект собирается и проходит все тесты.

---

## 📄 Лицензия

Проект распространяется под лицензией MIT.

---

**Swusher — арендуй и делись легко!**