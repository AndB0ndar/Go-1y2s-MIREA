# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY shared/ ./shared/
COPY services/graphql/go.mod services/graphql/go.sum* ./services/graphql/

WORKDIR /app/services/graphql
RUN go mod download

COPY services/graphql/ ./

RUN go build -o /graphql-server ./cmd/graphql


# Stage 2: Runner
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

COPY --from=builder /graphql-server .

EXPOSE 8090

CMD ["./graphql-server"]
