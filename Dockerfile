# Stage 1: Builder
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags "-s -w -extldflags '-static'" \
    -o /app/todo \
    ./cmd/main.go 

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup -H -s /sbin/nologin

WORKDIR /app

COPY --from=builder /app/todo .

RUN chown appuser:appgroup /app/todo && \
    chmod +x /app/todo

USER appuser

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=4s --start-period=10s --retries=4 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/me || exit 1  
CMD ["./todo"]