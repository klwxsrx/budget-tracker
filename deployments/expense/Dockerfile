# Download modules
FROM golang:1.14 as modules

ADD ./go.mod ./go.sum /
RUN go mod download


# Build the binary
FROM golang:1.14 as builder

RUN useradd -u 1001 appuser

COPY --from=modules /go/pkg /go/pkg
ADD . /build
WORKDIR /build

RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 \
    go build -o ./bin/expense ./cmd/expense

RUN chmod +x ./bin/expense


# Run the binary
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
USER appuser

COPY --from=builder /build/bin/expense /app/bin/expense

ADD ./data/mysql/migrations/eventstore /app/migrations/eventstore
ENV EVENT_STORE_MIGRATIONS_DIR=/app/migrations/eventstore

EXPOSE 80

CMD ["/app/bin/expense"]