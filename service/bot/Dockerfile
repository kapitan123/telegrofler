# syntax=docker/dockerfile:1

# ===================================================================================
# === Stage 1:Builder container =====================================================
# ===================================================================================
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Fetch modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy Go source files
# Copy Go source files
COPY common/ ./common/
COPY service/bot/ ./service/bot/

# All GCE instances are x86_64 linux based machines. Dynamic linking is disabled because it will be copied to scratch.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o server github.com/kapitan123/telegrofler/service/bot

# ===================================================================================
# === Stage 2: Create a lightweight container =======================================
# ===================================================================================
FROM scratch

# Add certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Add binary
COPY --from=builder /build/server /server

EXPOSE 9001

ENTRYPOINT [ "/server"]
