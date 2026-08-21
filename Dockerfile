# ---- Build stage ----
FROM golang:1.26.3-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Cache go mod downloads separately from source changes
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -o bin/seeder ./cmd/seeder && \
    CGO_ENABLED=0 GOOS=linux go build -o bin/keygen ./cmd/keygen

# ---- Run stage ----
FROM alpine:3.20
WORKDIR /app

# Reuse the cert bundle already fetched in the builder stage instead of
# running apk add again here — avoids a second hit to the Alpine mirror,
# which is where the earlier build failed (temporary CDN error).
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/bin/api .
COPY --from=builder /app/bin/seeder .
COPY --from=builder /app/bin/keygen .

EXPOSE 8080
# Default command is still the API server. seeder/keygen are run on demand
# via `docker compose run --entrypoint ./seeder api` (see `make docker-seed`).
ENTRYPOINT ["./api"]