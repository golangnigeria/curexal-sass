# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-exe --no-cache git ca-certificates tzdata

# Copy dependency manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o curexal-backend ./cmd/CUREXAL

# Final runtime stage
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# Copy binary & templates
COPY --from=builder /app/curexal-backend /app/curexal-backend
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/static /app/static
COPY --from=builder /app/database /app/database

EXPOSE 8080

CMD ["/app/curexal-backend"]
