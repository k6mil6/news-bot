FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY internal/storage/migrations /usr/local/src/news-bot/internal/storage/migrations

COPY ./ ./
RUN go build -o ./bin/app cmd/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /app
COPY --from=builder /usr/local/src/news-bot/internal/storage/migrations /migrations
COPY config.hcl /config.hcl

CMD ["/app"]