# Stage 1: Build
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o huifu-server ./cmd/server

# Stage 2: Runtime
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata curl && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

RUN addgroup -S huifu && adduser -S huifu -G huifu

COPY --from=builder /build/huifu-server /usr/local/bin/huifu-server
COPY --from=builder /build/config.yaml /etc/huifu/config.yaml
COPY --from=builder /build/migrations /etc/huifu/migrations

USER huifu
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/health || exit 1

ENTRYPOINT ["huifu-server"]
