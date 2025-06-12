FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine
RUN apk add --no-cache postgresql-client
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY run.sh /app/run.sh
RUN chmod +x /app/run.sh
WORKDIR /app

CMD ["/app/run.sh"]