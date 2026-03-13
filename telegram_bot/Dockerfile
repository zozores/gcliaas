# Stage 1: Build the Go binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o telegram-bot main.go

# Stage 2: Final lightweight image
FROM alpine:latest

WORKDIR /app

# Install ca-certificates required for HTTPS calls
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/telegram-bot .

# Run the application
CMD ["./telegram-bot"]
