# go_shortener

Сервис сокращения URL (HTTP + gRPC) с поддержкой пользовательских ссылок, batch-сокращений, расширения short->original, статистики и удаления.

## Стек
- Go 1.20
- HTTP: `chi`
- gRPC: `google.golang.org/grpc`
- Логирование: `zerolog`
- Auth: cookie `session_token` (JWT)
- Storage: PostgreSQL (pgx), file (JSON), in-memory
- TLS: генерация `cert.pem`/`key.pem` при включении

## Локальный запуск

### 1) SECRET_KEY
JWT подпись и валидация используют переменную окружения `SECRET_KEY`.

Пример:
```bash
export SECRET_KEY="change-me"
```

### 2) Сборка
```bash
cd cmd/shortener
go build -buildvcs=false -o shortener
```

### 3) Запуск
HTTP сервер слушает адрес из флага `-a` (по умолчанию `:8080`).

Базовый URL для формирования short-ссылок задаётся флагом `-b` (по умолчанию `http://localhost:8080`).

Запуск (по умолчанию пытается использовать PostgreSQL, т.к. `-d` имеет значение по умолчанию):
```bash
./shortener -a ":8080" -b "http://localhost:8080"
```

Если PostgreSQL не нужен и хотите file storage:
```bash
./shortener -d="" -f "/tmp/short-url-db.json"
```

Если PostgreSQL нужен (пример DSN):
```bash
./shortener -d="host=localhost user=shortener password=password dbname=shortener sslmode=disable"
```

TLS:
```bash
./shortener -s=true
```

### Trusted subnet (для `GET /api/internal/stats`)
Если задать `-t`/`TRUSTED_SUBNET`, то `GET /api/internal/stats` станет доступен по IP без cookie.
Реальный IP берётся из заголовка `X-Real-IP`.

Пример:
```bash
./shortener -t "127.0.0.1/32"
```

## HTTP API

### Создание короткой ссылки
`POST /`
- Body: plain text URL
- Ответ: `201` + plain text short URL

`POST /api/shorten`
- Body JSON: `{"url":"https://example.com"}`
- Ответ: `201` + JSON `{"result":"http://localhost:8080/<id>"}`

### Batch сокращение
`POST /api/shorten/batch`
- Body JSON (массив):
[
  {"correlation_id":"4444","original_url":"https://..."},
  {"correlation_id":"5555","original_url":"https://..."}
]
- Ответ `201`:
[
  {"correlation_id":"4444","short_url":"http://...","original_url":"https://..."},
  {"correlation_id":"5555","short_url":"http://...","original_url":"https://..."}
]

### Получение оригинала
`GET /{id}`
- 307 + `Location` (original URL)
- 204 если ссылки нет
- 410 если ссылка удалена

### Корзина пользователя
Требуется cookie `session_token`.

`GET /api/user/urls` -> `200` JSON bucket
`DELETE /api/user/urls` -> `202`, body: `["<id1>","<id2>",...]`

### Внутренняя статистика
Ответ:
`{"links_count":<int>,"users_count":<int>}`

Доступ:
- по trusted subnet (если задан `-t`)
- иначе — по cookie `session_token`

`GET /api/internal/stats`

### Health
`GET /ping` — `200` если storage доступно.

## gRPC API
Сервис `shortener.Shortener` (файл `internal/proto/shortener.proto`).

Методы:
- `Shorten(ShortRequest)` -> `ShortResponse`
- `Expand(ExpandRequest)` -> `ExpandResponse`
- `Batch(ShortBatchRequest)` -> `ShortBatchResponse`
- `GetUserBucket(GetUserBucketRequest)` -> `GetUserBucketResponse`
- `Stats(google.protobuf.Empty)` -> `StatResponse`
- `Delete(DeleteRequest)` -> `google.protobuf.Empty`

## Тестирование и CI
- Локально: `go test ./...`
- GitHub Actions запускает:
  - `go vet` через `statictest`
  - автотесты через `shortenertestbeta`
- Для автотестов используются ветки формата `iter<number>`.
