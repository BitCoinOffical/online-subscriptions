REST-сервис для управления онлайн подписками пользователей с поддержкой авторизации, фильтрации и агрегации данных.

## Стек

- **Go 1.26**
- **Gin** — HTTP фреймворк
- **PostgreSQL** — основная база данных
- **Redis** — хранение сессий и rate limiting
- **pgx/v5** — драйвер PostgreSQL
- **squirrel** — построитель SQL запросов
- **goose** — миграции
- **JWT** — авторизация
- **Swagger** — документация API
- **Docker / Docker Compose** — запуск сервисов
- **zap** — структурированное логирование
- **testcontainers** — интеграционные тесты
- **gomock** — юнит тесты

## Архитектура

Микросервисная архитектура с двумя независимыми сервисами:

## Структура проекта
```text
├── auth-service/           # Авторизация и управление сессиями
│   ├── cmd/
│   ├── config/
│   ├── internal/
│   │   ├── api/            # HTTP хендлеры, middleware
│   │   ├── domain/         # Модели, DTO, ошибки
│   │   └── interfaces/     # Сервисы, репозитории, кеш
│   ├── migrations/
│   └── pkg/                # JWT, логгер
│
└── subscription-service/   # Управление подписками
├── cmd/
├── config/
├── internal/
│   ├── api/            # HTTP хендлеры, middleware, тесты
│   ├── domain/         # Модели, DTO, ошибки
│   └── interfaces/     # Сервисы, репозитории, кеш, тесты
├── migrations/
└── pkg/                # Парсер дат, логгер

```
## API

### Auth Service (порт 8081)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/auth/register` | Регистрация |
| POST | `/auth/login` | Вход |
| POST | `/auth/refresh` | Обновление access токена |
| DELETE | `/auth/logout` | Выход |

### Subscription Service (порт 8080)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/api/v1/subscriptions` | Создать подписку |
| GET | `/api/v1/subscriptions/` | Получить все подписки |
| GET | `/api/v1/subscriptions/:id` | Получить подписку по ID |
| PATCH | `/api/v1/subscriptions/:id` | Частичное обновление |
| PUT | `/api/v1/subscriptions/:id` | Полное обновление |
| DELETE | `/api/v1/subscriptions/:id` | Удалить подписку |
| GET | `/api/v1/subscriptions` | Сумма подписок с фильтрами |

### Авторизация

Все запросы к subscription-service требуют заголовок: Authorization: Bearer <access_token>

### Фильтры для подсчёта суммы

| Параметр | Описание |
|----------|----------|
| `from` | Дата начала периода (MM-YYYY или DD-MM-YYYY) |
| `to` | Дата конца периода (MM-YYYY или DD-MM-YYYY) |
| `user_id` | UUID пользователя |
| `service_name` | Название сервиса |

## Запуск

1. Склонируй репозиторий:
```bash
git clone https://github.com/BitCoinOffical/online-subscriptions.git
cd online-subscriptions
```

2. Создай `.env` файлы на основе примеров:
```bash
cp auth-service/.env.example auth-service/.env
cp subscription-service/.env.example subscription-service/.env
```

3. Запусти через Docker Compose:
```bash
docker compose up --build
```

| Сервис | URL |
|--------|-----|
| Subscription API | http://localhost:8080 |
| Auth API | http://localhost:8081 |
| Swagger | http://localhost:8080/swagger/index.html |

## Тесты

```bash
go test ./...
```

Покрытие:
- **Хендлеры** — юнит тесты с gomock + httptest
- **Сервисы** — юнит тесты с gomock
- **Репозиторий** — интеграционные тесты с testcontainers (требуется Docker)

## Пример запроса

Регистрация:
```bash
curl -X POST http://localhost:8081/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "password_confirm": "password123"
  }'
```

Создание подписки:
```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```
