# ShareHub

ShareHub — это современная платформа для обмена и размещения объявлений, реализованная на Go (Echo), PostgreSQL и современном фронтенде.

---

## 🚀 Возможности

- Регистрация и вход по email/паролю
- Вход через Яндекс и Google (OAuth)
- Безопасное хранение паролей (bcrypt)
- JWT-аутентификация пользователей
- Валидация данных на клиенте и сервере
- Хранение профиля пользователя, объявлений, чатов и отзывов
- Получение координат пользователя при регистрации
- Современный адаптивный интерфейс

---

## 🛠️ Технологии

- **Backend:** Go, Echo, PostgreSQL, pgx, bcrypt, JWT
- **Frontend:** HTML5, CSS3, Vanilla JS (ES6+)
- **OAuth:** Яндекс, Google
- **Логирование:** zap
- **Безопасность:** bcrypt, JWT, CORS

---

## ⚡ Быстрый старт

### 1. Клонируйте репозиторий

```sh
git clone https://github.com/yourusername/sharehub.git
cd sharehub
```

### 2. Настройте переменные окружения

Создайте файл `.env` и укажите:

```
DATABASE_URL=postgres://user:password@localhost:5432/sharehub
JWT_SECRET=your_secret_key
YANDEX_CLIENT_ID=...
YANDEX_CLIENT_SECRET=...
GOOGLE_CLIENT_ID=...
```

### 3. Запустите PostgreSQL и создайте базу

```sh
createdb sharehub
```

### 4. Запустите сервер

```sh
go run ./cmd/server/main.go
```

### 5. Откройте в браузере

```
http://localhost:8080/
```

---

## 📂 Структура проекта

```
app/
  internal/
    handlers/      # HTTP-обработчики (регистрация, логин, OAuth)
    models/        # Модели данных (User, Advertisement и др.)
    static/        # Статические файлы (js, css)
    templates/     # HTML-шаблоны
  pkg/
    hash/          # Хэширование паролей (bcrypt)
    jwt/           # Работа с JWT-токенами
    store/         # Работа с базой данных (pgx)
```

---

## 🛡️ Безопасность

- Пароли хранятся только в виде bcrypt-хэшей
- JWT-токены используются для аутентификации
- Валидация данных на сервере и клиенте
- Защита от повторной регистрации по email

---

## 📖 Примеры API

### Регистрация

```http
POST /register/
Content-Type: application/json

{
  "name": "Иван",
  "lastname": "Иванов",
  "email": "ivan@example.com",
  "phone": "+79991234567",
  "password": "yourpassword"
}
```

### Вход

```http
POST /login/
Content-Type: application/json

{
  "email": "ivan@example.com",
  "password": "yourpassword"
}
```

---

## 📝 TODO

- Добавить восстановление пароля
- Улучшить дизайн профиля
- Реализовать push-уведомления
- Добавить фильтры и поиск по объявлениям

---

## 📄 Лицензия

MIT License

---

**ShareHub — делитесь и находите нужное легко!**