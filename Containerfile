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
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8088
CMD ["./main"]