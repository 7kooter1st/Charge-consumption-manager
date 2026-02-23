# Документация API

Спецификация OpenAPI хранится локально и отдаётся приложением:

- **Файл:** `internal/app/handler/openapi.json` (встроен через `//go:embed`)
- **Спецификация:** `GET /swagger/doc.json`
- **Swagger UI:** `GET /swagger/` или `GET /swagger/index.html`

По умолчанию сервер слушает порт **8000** (из `config.toml`). Чтобы запустить на порту 8080, задайте в `.env`:
`SERVICE_PORT=8080`.

Откройте в браузере: `http://localhost:8000/swagger/` (или `http://localhost:8080/swagger/`, если указали порт 8080).
