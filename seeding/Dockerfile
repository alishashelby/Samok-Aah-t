FROM golang:1.22-alpine AS builder

WORKDIR /faker
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o seeding ./cmd

FROM alpine:latest
WORKDIR /faker
COPY --from=builder /faker/seeding .
CMD ["./seeding"]