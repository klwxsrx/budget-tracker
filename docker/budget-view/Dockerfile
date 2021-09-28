# Create user
FROM alpine:latest AS builder

RUN adduser --disabled-password --uid=1001 appuser

# Run the binary
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
USER appuser

COPY ./bin/budget-view /app/bin/budget

EXPOSE 8080

CMD ["/app/bin/budget"]