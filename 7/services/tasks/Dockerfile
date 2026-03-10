# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum* ./
COPY services/tasks/go.mod services/tasks/go.sum* ./services/tasks/
COPY services/auth/go.mod services/auth/go.sum* ./services/auth/
COPY shared/go.mod shared/go.sum* ./shared/

RUN go mod download

COPY services/tasks/ ./services/tasks/
COPY shared/ ./shared/

WORKDIR /app/services/tasks

RUN go build -o /tasks-server ./cmd/tasks


# Stage 2: Runner
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

COPY --from=builder /tasks-server .

EXPOSE 8082

CMD ["./tasks-server"]

