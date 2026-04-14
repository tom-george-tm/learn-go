# ─────────────────────────────────────────────────────────────────────────────
# Stage 1 — Builder
# Compiles the Go source into a single binary.
# ─────────────────────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files first so Docker caches this layer separately.
# "go mod download" only re-runs when go.mod or go.sum change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 and GOOS=linux produce a self-contained Linux binary.
# -s -w strips debug info to reduce binary size.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/server \
    ./cmd/main.go


# ─────────────────────────────────────────────────────────────────────────────
# Stage 2 — Final Image
# Copies only the binary into a minimal Alpine image (~5MB vs ~800MB).
# ─────────────────────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

# Run as a non-root user for security.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server .

USER appuser

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:3000/health || exit 1

CMD ["./server"]