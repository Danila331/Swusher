# 🚀 ShareHub Backend & Infrastructure

> **ShareHub** — платформа p2p-аренды вещей между частными лицами. Это технический репозиторий backend-сервиса, API, DevOps-инфраструктуры и сопутствующих компонентов.

---

## 🧱 Архитектура

Микросервисная архитектура с модулями:

- `users` — регистрация, авторизация, профили
- `ads` — объявления, категории, модерация
- `rents` — бронирование, статусы, история аренд
- `payments` — интеграция с платёжными системами (например, ЮKassa, Tinkoff)
- `notifications` — email / push / in-app уведомления
- `search` — полнотекстовый поиск (Elastic)
- `media` — загрузка и хранение файлов (S3-compatible)
- `ml-core` — рекомендательная система и динамическое ценообразование

---

## ⚙️ Технологии

| Слой | Технологии |
|------|------------|
| Язык | Go 1.22 / Python 3.11 (ML-сервисы) |
| API | gRPC / REST (OpenAPI 3) |
| Хранение | PostgreSQL / Redis / S3 |
| Очереди | Kafka |
| Кеширование | Redis |
| CI/CD | GitHub Actions + Docker + Kubernetes |
| Мониторинг | Prometheus + Grafana |
| Тесты | Unit / Integration / E2E (Testcontainers) |

---

## 🚀 Быстрый старт (Dev)

### 🖥 Требования

- Docker + Docker Compose
- Go ≥ 1.22
- PostgreSQL ≥ 14 (если без docker)
- Make

### 🔧 Запуск локально

```bash
git clone https://github.com/sharehub-team/platform-backend.git
cd platform-backend

# Инициализация env и запуск контейнеров
cp .env.example .env
make dev-up
````

Платформа поднимается на:

* API: `http://localhost:8080/`
* Swagger: `http://localhost:8080/docs`
* PGAdmin: `http://localhost:8081`
* Mailhog: `http://localhost:8025`

---

## 📦 Структура репозитория

```
.
├── cmd/                # entrypoints
├── internal/
│   ├── handlers/
│   ├── midlewary/
│   ├── models/
|   |   |── advertisements/
|   |   |── chats/
|   |   |── passports/
|   |   |── reviews/
|   |   |── users/
│   |── servers/
|   |── static/
|   |   |── css/
|   |   |── js/
|   |   |── images/
|   └── templates/         
├── pkg/                # переиспользуемые библиотеки
├── api/                # OpenAPI / protobuf схемы
├── migrations/         # SQL миграции
├── scripts/            # служебные скрипты
├── deploy/             # Helm-чарты и K8s манифесты
└── docs/               # документация по архитектуре
```

---

## 📈 Мониторинг и логирование

* **Prometheus** собирает метрики (latency, throughput, errors)
* **Grafana**: дашборды (пользователи, аренды, конверсии)
* **Loki** + **Grafana**: централизованные логи
* **Jaeger**: distributed tracing

---

## 🧪 Тестирование

```bash
make test         # Unit
make test-int     # Интеграционные
make test-e2e     # Сквозные (e2e)
```

---

## 🔐 Безопасность

* OAuth2 / JWT авторизация
* RBAC по ролям (пользователь / владелец / модератор / админ)
* Rate limiting + защита от повторных запросов
* GDPR / cookie consent / secure headers

---

## 🧠 ML-модули

Интеграция с модулем рекомендаций:

* Автоматические рекомендации похожих объявлений
* Расчёт вероятности аренды и оптимизация цены
* Фреймворк: **FastAPI + LightGBM + PostgreSQL**

---

## 📦 Deployment (CI/CD)

* CI: GitHub Actions
* CD: Helm + Kubernetes (Dev, Staging, Prod)
* Preview-окружения для Pull Request-ов
* Автокаталогизация через ArgoCD

---

## 💬 Документация

* Swagger: `/docs`
* API спецификации: `./api/openapi.yaml`
* Архитектура: `./docs/architecture.md`
* Dev-инструкции: `./docs/contributing.md`

---

## 📫 Контакты

* CTO: `tech@sharehub.ru`
* Telegram: [@sharehub\_dev](https://t.me/sharehub_dev)
* Поддержка: [support@sharehub.ru](mailto:support@sharehub.ru)

---

## 📝 Лицензия

Проект распространяется под лицензией MIT.