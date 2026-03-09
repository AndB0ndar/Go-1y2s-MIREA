FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ./services/tasks/go.* ./
RUN go mod download

COPY . .

RUN go build -o /tasks-app ./services/tasks/cmd/tasks

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /tasks-app .
EXPOSE 8082
CMD ["./tasks-app"]
