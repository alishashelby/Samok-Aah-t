FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git postgresql-client
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /test
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o migration_tester ./cmd

FROM alpine
RUN apk add --no-cache postgresql-client
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /test/migration_tester /app/migration_tester

WORKDIR /app
CMD ["/app/migration_tester"]