# Online Subscriptions

REST-сервис для агрегации данных об онлайн подписках пользователей.

Тестовое задание Junior Golang Developer — Effective Mobile.

## Стек

- **Go 1.26**
- **Gin** — HTTP фреймворк
- **PostgreSQL** — база данных
- **pgx/v5** — драйвер PostgreSQL
- **squirrel** — построитель SQL запросов
- **golang-migrate** — миграции
- **Swagger** — документация API
- **Docker / Docker Compose** — запуск сервиса
- **zap** — логирование
- **testcontainers** — интеграционные тесты
- **gomock** — юнит тесты

## Структура проекта
```text
├── cmd/                    # Точка входа
├── config/                 # Конфигурация
├── docs/                   # Swagger документация
├── migrations/             # SQL миграции
├── pkg/
│   ├── date/               # Парсер дат
│   └── logger/             # Логгер
└── internal/
├── api/
│   ├── handlers/       # HTTP хендлеры + тесты
│   └── response/       # Хелперы ответов
├── domain/
│   ├── dto/            # DTO
│   ├── models/         # Модели
│   └── errors.go       # Доменные ошибки
└── interfaces/http/
├── repo/           # Репозиторий + тесты
└── services/       # Сервисы + тесты
```
## API

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/api/v1/subscriptions` | Создать подписку |
| GET | `/api/v1/subscriptions/` | Получить все подписки |
| GET | `/api/v1/subscriptions/:id` | Получить подписку по ID |
| PATCH | `/api/v1/subscriptions/:id` | Частичное обновление |
| PUT | `/api/v1/subscriptions/:id` | Полное обновление |
| DELETE | `/api/v1/subscriptions/:id` | Удалить подписку |
| GET | `/api/v1/subscriptions` | Сумма подписок с фильтрами |

### Фильтры для подсчёта суммы

| Параметр | Описание |
|----------|----------|
| `from` | Дата начала периода (MM-YYYY) или (DD-MM-YYYY)|
| `to` | Дата конца периода (MM-YYYY) или (DD-MM-YYYY)|
| `user_id` | UUID пользователя |
| `service_name` | Название сервиса |

### Пример запроса

```json
POST /api/v1/subscriptions
{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025",
    "end_date": "12-2025"
}
```

## Запуск

1. Склонируй репозиторий
2. Создай `.env` файл на основе `.env.example`
3. Запусти через Docker Compose:

```bash
make build или docker compose up --build
```

Сервис будет доступен на `http://localhost:8080`

Swagger UI: `http://localhost:8080/swagger/index.html`

## Тесты

```bash
make test или go test ./...
```

Покрытие:
- **Хендлеры** — юнит тесты с gomock + httptest
- **Сервисы** — юнит тесты с gomock
- **Репозиторий** — интеграционные тесты с testcontainers (требуется Docker)