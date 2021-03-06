# Download modules
FROM golang:1.16 AS modules

COPY ./go.mod ./go.sum /
RUN go mod download


# Build the binary
FROM golang:1.16 AS builder

RUN useradd -u 1001 appuser

COPY --from=modules /go/pkg /go/pkg
COPY . /build
WORKDIR /build

RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 \
    go build -o ./bin/budget ./cmd/budget

RUN chmod +x ./bin/budget


# Run the binary
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
USER appuser

COPY --from=builder /build/bin/budget /app/bin/budget

COPY --from=builder /build/data/mysql/migrations/budget /app/migrations
ENV DATABASE_MIGRATIONS_DIR=/app/migrations

EXPOSE 8080

CMD ["/app/bin/budget"]