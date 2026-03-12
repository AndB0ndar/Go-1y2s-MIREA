# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY shared/ ./shared/
COPY services/worker/go.mod services/worker/go.sum* ./services/worker/

WORKDIR /app/services/worker
RUN go mod download

COPY services/worker/ ./

RUN go build -o /worker ./cmd/worker


# Stage 2: Runner
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

COPY --from=builder /worker .

CMD ["./worker"]
