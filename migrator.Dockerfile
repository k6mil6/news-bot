FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY internal/storage/migrations /usr/local/src/news-bot/internal/storage/migrations

COPY ./ ./
RUN go build -o ./bin/app cmd/migrator/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /app
COPY --from=builder /usr/local/src/news-bot/internal/storage/migrations /migrations
COPY config.hcl /config.hcl

CMD ["/app"]