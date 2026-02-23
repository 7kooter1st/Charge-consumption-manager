# Запуск с Docker (docker-compose)

## Контейнеры

| Сервис          | Порт  | Описание                          |
|-----------------|-------|-----------------------------------|
| postgres        | 5432  | БД `consumption_db`               |
| redis           | 6379  | Redis (пароль: `redispassword`)   |
| swagger-editor  | 8080  | Редактор OpenAPI                  |
| minio           | 9000, 9001 | Хранилище (консоль :9001)   |
| adminer         | 8081  | Веб-интерфейс к PostgreSQL        |

Приложение (Go) запускается на хосте и слушает порт **8000** (config.toml).

## Перед первым запуском

1. Скопируйте переменные под контейнеры:
   ```bash
   cp .env.example .env
   ```
2. Поднимите контейнеры:
   ```bash
   docker compose up -d
   ```
3. Примените миграции и при необходимости заполните данные:
   ```bash
   go run cmd/migrate/main.go
   # при необходимости: psql -h localhost -U postgres -d consumption_db -f sql.sql
   ```
4. Запустите приложение:
   ```bash
   go run cmd/start/main.go
   ```

API: `http://127.0.0.1:8000`  
Swagger Editor: `http://localhost:8080`  
Adminer: `http://localhost:8081`  
MinIO Console: `http://localhost:9001`
