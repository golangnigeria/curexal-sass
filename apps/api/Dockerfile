# Build stage
FROM golang:alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto
ENV CGO_ENABLED=0

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy dependency manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN GOOS=linux go build -ldflags="-w -s" -o bin/curexal-backend ./cmd/CUREXAL

# Final runtime stage
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# Copy binary & assets
COPY --from=builder /app/bin/curexal-backend /app/curexal-backend
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/static /app/static
COPY --from=builder /app/database /app/database

EXPOSE 8080

CMD ["/app/curexal-backend"]
