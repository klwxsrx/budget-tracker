# Build the binary
FROM golang:1.15 as builder

RUN useradd -u 1001 appuser

COPY . /build
WORKDIR /build

RUN go mod download

RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 \
    go build -o ./bin/expense ./cmd/expense

RUN chmod +x ./bin/expense


# Run the binary
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
USER appuser

COPY --from=builder /build/bin/expense /app/bin/expense

COPY --from=builder /build/data/mysql/migrations/expense /app/migrations
ENV MIGRATIONS_DIR=/app/migrations

EXPOSE 8080

CMD ["/app/bin/expense"]