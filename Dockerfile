# ---- Build stage ----
FROM golang:1.26.3-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

# Cache go mod downloads separately from source changes
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/api ./cmd/api

# ---- Run stage ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/bin/api .

EXPOSE 8080
ENTRYPOINT ["./api"]