FROM golang:1.20-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/gochat ./cmd/gochat

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata wget grep

WORKDIR /app

COPY --from=builder /out/gochat /app/gochat
COPY configs /app/configs

RUN mkdir -p /app/static/avatars /app/static/files /app/logs

EXPOSE 8000

CMD ["/app/gochat"]
