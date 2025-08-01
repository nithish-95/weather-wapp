FROM golang:alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for better cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o /app/bin/weatherapp cmd/server/main.go

FROM alpine:latest

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary and required files
COPY --from=builder /app/bin/weatherapp .
COPY --from=builder /app/templates ./templates

# Set permissions
RUN chown -R appuser:appgroup /app
USER appuser

EXPOSE 3000

CMD ["./weatherapp"]

