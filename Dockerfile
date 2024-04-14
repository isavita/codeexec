# Build stage
FROM golang:1.22 AS build

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the compiled binary from the build stage
COPY --from=build /app/main .

# Set the default value for the PORT environment variable
ENV PORT=8080

# Expose the port on which the server will run
EXPOSE ${PORT}

# Run the application
CMD ["./main"]
