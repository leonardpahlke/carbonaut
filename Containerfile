FROM golang:1.22-alpine AS builder
RUN apk --no-cache add build-base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main .

FROM alpine:latest
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/main .
RUN echo -e '#!/bin/sh\n\
# Check if CONFIG_PATH is set\n\
if [ -z "$CONFIG_PATH" ]; then\n\
  echo "Error: CONFIG_PATH is not set."\n\
  exit 1\n\
fi\n\
# Start the application with the specified configuration file\n\
exec ./main -c "$CONFIG_PATH"\n' > /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8088
ENTRYPOINT ["/app/entrypoint.sh"]
