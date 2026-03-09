FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ./services/auth/go.mod ./
RUN go mod download

COPY . .

RUN go build -o /auth-app ./services/auth/cmd/auth

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /auth-app .
EXPOSE 8081
CMD ["./auth-app"]
