# Практическое задание 7. Написание Dockerfile и сборка контейнера

**Студент:** Бондарь Андрей Ренатович  
**Группа:** ЭФМО-02-25

---

## Цель работы
Научиться упаковывать сервис в Docker-образ с использованием multi-stage сборки, обеспечивая воспроизводимость и компактность образа, а также запускать связанные сервисы через docker-compose.

---

## Dockerfile для сервисов

### Общая структура (multi-stage build)
Для обоих сервисов (`auth` и `tasks`) используется одинаковый подход:
- **Builder stage** – компиляция бинарного файла в полном образе Go.
- **Runner stage** – минимальный образ Alpine, содержащий только бинарник, необходимые системные зависимости (ca-certificates) и curl для healthcheck в docker-compose.

### Dockerfile для сервиса `auth`
**Файл:** `services/auth/Dockerfile`
```dockerfile
# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей (для кеширования)
COPY services/auth/go.mod services/auth/go.sum ./services/auth/
COPY shared/go.mod shared/go.sum ./shared/

WORKDIR /app/services/auth
RUN go mod download

# Копируем остальной код
COPY services/auth/ ./services/auth/
COPY shared/ ./shared/

# Сборка бинарного файла
RUN go build -o /auth-server ./cmd/auth

# Stage 2: Runner
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

COPY --from=builder /auth-server .

EXPOSE 8081 50051

CMD ["./auth-server"]
```

### Dockerfile для сервиса `tasks`

Файл `services/tasks/Dockerfile` аналогично `auth`.

### `.dockerignore` для каждого сервиса
Для исключения лишних файлов из контекста сборки созданы файлы `.dockerignore`:

**`services/auth/.dockerignore`**
```
.git
.gitignore
README.md
tmp/
bin/
*.log
.dockerignore
Dockerfile
```

**`services/tasks/.dockerignore`** – аналогичный.

---

## Переменные окружения сервисов

Сервисы получают конфигурацию через переменные окружения, которые передаются при запуске контейнера.

| Сервис | Переменная       | Значение по умолчанию     | Описание                          |
|--------|------------------|---------------------------|-----------------------------------|
| auth   | `AUTH_PORT`      | 8081                      | HTTP-порт для Auth                |
| auth   | `AUTH_GRPC_PORT` | 50051                     | gRPC-порт для Auth                |
| tasks  | `TASKS_PORT`     | 8082                      | HTTP-порт для Tasks                |
| tasks  | `AUTH_GRPC_ADDR` | auth:50051                | Адрес gRPC-сервера Auth            |
| tasks  | `DB_HOST`        | postgres                  | Хост PostgreSQL                    |
| tasks  | `DB_PORT`        | 5432                      | Порт PostgreSQL                    |
| tasks  | `DB_USER`        | tasks_user                | Пользователь БД                    |
| tasks  | `DB_PASSWORD`    | tasks_pass                | Пароль БД                          |
| tasks  | `DB_NAME`        | tasks_db                  | Имя БД                             |

Все переменные задаются в `docker-compose.yml` в секции `environment`.

---

## Запуск через docker-compose (взаимодействие в одной сети)

Файл `deploy/docker-compose.yml` объединяет все сервисы: PostgreSQL, Auth, Tasks и NGINX (для HTTPS). Сервисы общаются внутри общей сети `app-network` по именам контейнеров.

**Фрагмент `docker-compose.yml`:**
```yaml
include:  #  для контейнеров nginx и мониторинга
  - monitoring/docker-compose.monitoring.yml
  - tls/docker-compose.tls.yml

services:
  postgres:
    image: postgres:15
    # ...

  auth:
    build:
      context: ..
      dockerfile: services/auth/Dockerfile
    environment:
      - AUTH_PORT=8081
      - AUTH_GRPC_PORT=50051
    networks:
      - app-network

  tasks:
    build:
      context: ..
      dockerfile: services/tasks/Dockerfile
    environment:
      - TASKS_PORT=8082
      - AUTH_GRPC_ADDR=auth:50051
      - DB_HOST=postgres
      # ...
    depends_on:
      - auth
      - postgres
    networks:
      - app-network
```

**Важно:** Внутри docker-сети Tasks обращается к Auth по имени `auth:50051`, а не через `localhost`. Это стандартный способ взаимодействия контейнеров.

---

## Команды сборки и запуска

### Сборка образов вручную (из корня проекта)
```bash
# Auth
cd services/auth
docker build -t techip-auth:latest .

# Tasks
cd ../tasks
docker build -t techip-tasks:latest .
```

### Запуск через docker-compose (рекомендуется)
```bash
cd deploy
docker-compose up -d --build
```

### Просмотр логов
```bash
docker-compose logs -f
```

### Остановка и удаление контейнеров
```bash
docker-compose down
```

---

## Проверка работоспособности (curl)

Работоспособность можно посмотреть в предыдущих практиках.

---

## 8. Выводы
- Для обоих сервисов написаны multi-stage Dockerfile, обеспечивающие минимальный размер образов.
- Настроен запуск через docker-compose с общей сетью, что гарантирует надёжное взаимодействие контейнеров.
- Конфигурация передаётся через переменные окружения, секреты не вшиты в образы.
- Проведена успешная проверка работы всех эндпоинтов через curl.

