# Create user
FROM alpine:latest AS builder

RUN adduser --disabled-password --uid=1001 appuser

# Run the binary
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
USER appuser

COPY ./bin/budget /app/bin/budget

COPY ./data/mysql/migrations/budget /app/migrations
ENV DATABASE_MIGRATIONS_DIR=/app/migrations

EXPOSE 8080

CMD ["/app/bin/budget"]