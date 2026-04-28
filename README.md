# Task Service

Сервис для управления задачами с HTTP API на Go.

## Требования

- Go `1.23+`
- Docker и Docker Compose

## Быстрый запуск через Docker Compose

```bash
docker compose up --build
```

После запуска сервис будет доступен по адресу `http://localhost:8080`.

Если `postgres` уже запускался ранее со старой схемой, пересоздай volume:

```bash
docker compose down -v
docker compose up --build
```

Причина в том, что SQL-файл из `migrations/0001_create_tasks.up.sql` монтируется в `docker-entrypoint-initdb.d` и применяется только при инициализации пустого data volume.

## Swagger

Swagger UI:

```text
http://localhost:8080/swagger/
```

OpenAPI JSON:

```text
http://localhost:8080/swagger/openapi.json
```

## API

Базовый префикс API:

```text
/api/v1
```

Основные маршруты:

- `POST /api/v1/tasks`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/{id}`
- `PUT /api/v1/tasks/{id}`
- `DELETE /api/v1/tasks/{id}`

### Периодичность задач

Для создания и обновления задачи можно передавать поля:

- `recurrence_kind`: `none`, `daily`, `weekly`, `monthly_dates`, `monthly_even_days`, `monthly_odd_days`
- `recurrence_days`: массив чисел

Правила валидации:

- `none`, `daily`, `monthly_even_days`, `monthly_odd_days` -> `recurrence_days` должен быть пустым
- `weekly` -> `recurrence_days` обязателен, значения от `1` до `7` (дни недели)
- `monthly_dates` -> `recurrence_days` обязателен, значения от `1` до `31` (числа месяца)
