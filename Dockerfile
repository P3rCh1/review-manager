FROM golang:1.24.10 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o review-manager ./cmd/review-manager/main.go

RUN CGO_ENABLED=0 go build -o migrate ./cmd/migrate/main.go

FROM alpine:latest

COPY --from=builder /app/review-manager .
COPY --from=builder /app/migrate ./
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/config.yaml ./

EXPOSE 8080

ENTRYPOINT ["./review-manager", "-c", "config.yaml"]
