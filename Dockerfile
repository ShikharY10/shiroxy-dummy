# Use an official Golang runtime as the base image
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o crud-api-server

# Set environment variables for Gin mode and server port
ENV GIN_MODE=release
ENV PORT=8080

# Expose the port the server will listen on
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./crud-api-server"]
