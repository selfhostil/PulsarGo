# Use the official Golang image as the base image
FROM golang:1.18-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod file (if you're using Go modules)
COPY go.mod ./

# Optionally download Go modules (if you are using them)
RUN go mod download || true

# Copy the rest of the application code to the working directory
COPY . .

# Build the Go app
RUN go build -o pulsar .

# Expose port 5001 to the outside world
EXPOSE 5001

# Command to run the application
CMD ["./pulsar"]
