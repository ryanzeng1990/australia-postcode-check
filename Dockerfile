# 1. Build Stage: Compile the Go application
FROM golang:1.24.3-alpine AS builder

# Set necessary environment variables
ENV CGO_ENABLED=0
ENV GOOS=linux

# Install git and other build dependencies needed for goquery
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy the source code
COPY postcode_scraper.go .
COPY go.mod .
COPY go.sum .

# Download dependencies (this is needed because we use goquery)
RUN go mod download

# Build the application
# We use -o api to name the final executable 'api'
RUN go build -ldflags "-s -w" -o api postcode_scraper.go

# 2. Final Stage: Create a minimal production image
FROM alpine:latest

# Expose the port the application runs on
EXPOSE 8080

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/api .

# Run the executable
CMD ["./api"]