# Stage 1: Build the Go binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
# Copy go.sum if it exists (ignoring error if not)
COPY go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o gemini-listen main.go

# Stage 2: Final image with Homebrew and Gemini CLI
FROM homebrew/brew:latest

WORKDIR /app

# Install Gemini CLI via Homebrew
RUN brew install gemini-cli && brew cleanup && rm -rf /home/linuxbrew/.cache

# Copy the binary from the builder stage
COPY --from=builder /app/gemini-listen .

# Install Gemini Extensions
# RUN gemini extensions install https://github.com/github/github-mcp-server

# Expose the port
EXPOSE 8765

# Run the application
CMD ["./gemini-listen"]
