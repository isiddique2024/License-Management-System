# Stage 1: Build
FROM golang:1.22.5-alpine as builder

WORKDIR /usr/src/app

# Install dependencies needed for building, including gcc
RUN apk update && apk add --no-cache gcc musl-dev sqlite-dev git curl

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project directory except those ignored in .dockerignore
COPY . .

# Install Swag for documentation generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
RUN swag init -g main.go --output ./docs

# Build the Go application
RUN go build -o myapp

# Install Air for live reloading
RUN go install github.com/air-verse/air@v1.52.3

# Stage 2: Runtime
FROM golang:1.22.5-alpine

WORKDIR /app

# Install runtime dependencies, including gcc if using CGO
RUN apk add --no-cache gcc musl-dev curl

# Copy the source code and configuration files
COPY --from=builder /usr/src/app /app

# Copy Air binary from Go bin directory
COPY --from=builder /go/bin/air /usr/bin/air

# Expose the application port
EXPOSE 8001

# Set the entrypoint to use Air for live reloading
CMD ["air", "-c", ".air.toml"]
