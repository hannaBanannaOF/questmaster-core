# Stage 1: Builder
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to leverage Docker's build cache
COPY go.mod .
COPY go.sum .

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application binary
# CGO_ENABLED=0 disables CGo, resulting in a statically linked binary without external runtime dependencies
# GOOS=linux and GOARCH=amd64 target the Linux operating system and AMD64 architecture
# -o specifies the output filename for the executable
# -ldflags="-s -w" removes debugging information and symbol tables, reducing binary size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o questmaster-core -ldflags="-s -w" /app/cmd/app/main.go

# Stage 2: Runtime
# Use a minimal base image like scratch or a distroless image
FROM scratch

# Copy the built binary from the builder stage
COPY --from=builder /app/questmaster-core /questmaster-core

# Expose any ports your application listens on (optional)
EXPOSE 8080

# Specify the command to run when the container starts
ENTRYPOINT ["/questmaster-core"]