# Use the official Go image from Docker Hub
FROM golang:1.18-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules and Sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod tidy

# Copy the rest of the application files
COPY . .

# Build the Go application
RUN go build -o main ./cmd

# Expose port 8080
EXPOSE 8080

# Run the Go application
CMD ["./main"]
