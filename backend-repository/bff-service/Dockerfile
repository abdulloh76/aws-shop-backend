# Start from the official Golang image for building the application
FROM golang:1.20-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /bff-service

# Start a new stage from scratch for a smaller image
FROM alpine:latest  

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /bff-service /bff-service

# Copy the .env file into the container
COPY .env /app/.env

# Set the environment variables
ENV PORT=8080

# Expose the service port
EXPOSE 8080

# Command to run the executable
CMD ["/bff-service"]