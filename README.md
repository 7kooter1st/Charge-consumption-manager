# Charge Consumption Manager - Сервер для управления заявками на потребление электроэнергии

## Описание проекта

Веб-сервис для управления заявками на потребление электроэнергии. Поддерживает:
- Управление сценариями использования (use cases)
- Создание и управление заявками на потребление
- Модерацию заявок
- Хранение данных в PostgreSQL
- Хранение изображений в MinIO

## Технологический стек

- **Язык:** Go 1.25+
- **Фреймворк:** Gin Web Framework
- **ORM:** GORM
- **База данных:** PostgreSQL
- **Хранилище файлов:** MinIO
- **Конфигурация:** Viper

## Структура проекта

```
Charge-consumption-manager/
├── cmd/
│   ├── start/          # Точка входа приложения
│   └── migrate/        # Миграции БД
├── internal/
│   └── app/
│       ├── config/     # Конфигурация
│       ├── ds/         # Модели данных
│       ├── DTO/        # Data Transfer Objects
│       ├── dsn/        # Настройки подключения к БД
│       ├── handler/    # HTTP обработчики
│       ├── repository/ # Слой работы с БД
│       └── service/    # Бизнес-логика
├── config/
│   └── config.toml     # Файл конфигурации
├── sql.sql             # Тестовые данные
├── .env                # Переменные окружения
└── go.mod
```

## Установка и запуск

### Предварительные требования

1. Go 1.25 или выше
2. Docker и Docker Compose
3. PostgreSQL (через Docker)
4. MinIO (через Docker)

### 1. Настройка переменных окружения

Создайте файл `.env` в корне проекта:

```env
# PostgreSQL
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=charge_consumption_db
POSTGRES_HOST=localhost
POSTGRES_PORT=5432

# MinIO
MINIO_HOST=localhost
MINIO_PORT=9000
MINIO_ACCESS_KEY=minio
MINIO_SECRET_KEY=minio124
MINIO_BUCKET=usecase-images

# Application
SERVER_PORT=8080
```

### 2. Запуск PostgreSQL и MinIO

```bash
# Запустите PostgreSQL
docker run -d \
  --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=charge_consumption_db \
  -p 5432:5432 \
  postgres:latest

# Запустите MinIO
docker run -d \
  --name minio \
  -e MINIO_ROOT_USER=minio \
  -e MINIO_ROOT_PASSWORD=minio124 \
  -p 9000:9000 \
  -p 9001:9001 \
  minio/minio server /data --console-address ":9001"
```

### 3. Применение миграций

```bash
# Войдите в контейнер PostgreSQL
docker exec -it postgres psql -U postgres -d charge_consumption_db

# Выполните SQL-скрипты создания таблиц и заполнения тестовыми данными
\i /path/to/sql.sql
```

Или выполните миграции через приложение:

```bash
go run cmd/migrate/main.go
```

### 4. Сборка и запуск приложения

```bash
# Установите зависимости
go mod download

# Соберите приложение
go build -o main cmd/start/main.go

# Запустите приложение
./main
```

Или запустите напрямую:

```bash
go run cmd/start/main.go
```

Сервер будет доступен по адресу: `http://localhost:8080`

## API Endpoints

### Базовый URL
Все эндпоинты начинаются с `/api`

### Основные маршруты:

**Сценарии использования:**
- `GET /api/usecases/` - список сценариев
- `GET /api/usecases/:id` - один сценарий
- `POST /api/usecases/` - создать сценарий
- `PUT /api/usecases/:id` - обновить сценарий
- `DELETE /api/usecases/:id` - удалить сценарий
- `PUT /api/usecases/:id/image` - добавить изображение

**Заявки:**
- `GET /api/consumptions/` - список заявок (с фильтрацией)
- `GET /api/consumptions/draft` - черновик текущего пользователя
- `GET /api/consumptions/:id` - одна заявка
- `POST /api/consumptions/` - создать заявку
- `PUT /api/consumptions/:id/formate` - сформировать заявку
- `PUT /api/consumptions/:id/moderate` - модерация заявки
- `DELETE /api/consumptions/:id` - удалить черновик

**Управление сценариями в заявке:**
- `POST /api/consumptions/:id/usecases` - добавить сценарий
- `PUT /api/consumptions/:id/usecases/:usecase_id` - изменить длительность
- `DELETE /api/consumptions/:id/usecases/:usecase_id` - удалить сценарий

**Пользователи:**
- `POST /api/users/register` - регистрация
- `GET /api/users/:id` - данные пользователя
- `PUT /api/users/:id` - обновить данные

Подробные примеры см. в файле [API_EXAMPLES.md](../API_EXAMPLES.md)

## Важные изменения в версии 3.0

### Фиксированный пользователь
В текущей версии (без авторизации) используется фиксированный пользователь через singleton-паттерн:
- ID пользователя: 1 (admin)
- Параметр `?user_id=X` больше НЕ ТРЕБУЕТСЯ в URL

### Префикс /api
Все эндпоинты теперь доступны с префиксом `/api`:
- Старый: `/consumptions`
- Новый: `/api/consumptions`

## Бизнес-логика

### Статусы заявок
- `черновик` - создана, редактируется создателем
- `сформирован` - готова к модерации
- `завершен` - одобрена модератором
- `отклонен` - отклонена модератором
- `удален` - помечена как удаленная

### Правила изменения статусов
- **Создатель:**
  - Может удалять и формировать только черновики
  - При формировании рассчитывается общее потребление
  
- **Модератор:**
  - Может отклонять и завершать только сформированные заявки
  - При завершении/отклонении проставляется модератор и дата модерации

### Системные поля
Следующие поля рассчитываются на бэкенде и НЕ передаются с клиента:
- ID записей
- Статусы
- Создатель и модератор
- Даты создания, формирования, модерации

## Разработка

### Запуск в режиме разработки

```bash
# Включите автоперезагрузку с помощью air
go install github.com/cosmtrek/air@latest
air
```

### Проверка кода

```bash
# Форматирование
go fmt ./...

# Линтер
golangci-lint run

# Тесты
go test ./...
```

## Тестовые данные

В файле `sql.sql` содержатся тестовые данные:
- 5 пользователей (включая модераторов)
- 5 сценариев использования
- 5 заявок с различными статусами
- Связи между заявками и сценариями

## Структура базы данных

### Таблицы:
- `users` - пользователи системы
- `use_cases` - сценарии использования
- `consumptions` - заявки на потребление
- `usecase_consumptions` - связь многие-ко-многим (заявка ↔ сценарий)

## Известные ограничения

1. Авторизация - заглушка (используется фиксированный пользователь)
2. Аутентификация будет реализована в лабораторной работе №4
3. Загрузка файлов в MinIO требует дополнительной настройки CORS

## Лицензия

Учебный проект - Лабораторная работа №3

## Контакты

Если возникли вопросы по проекту, обратитесь к преподавателю или на GitHub Issues.
