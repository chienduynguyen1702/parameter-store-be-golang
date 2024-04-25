# Use the official Golang image as the base image
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -v -o ./main

# Start a new stage from scratch
FROM alpine:latest  

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage to the current stage
COPY --from=builder /app/main ./main

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/root/main"]