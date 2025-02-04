# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o url-shortener ./cmd/api

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/url-shortener .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./url-shortener"]