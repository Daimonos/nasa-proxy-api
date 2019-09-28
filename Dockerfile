# Dockerfile References: https://docs.docker.com/engine/reference/builder/
# Start from the latest golang base image
FROM golang:latest
# Add Maintainer Info
LABEL maintainer="James Hughes <jim.d.hughes@gmail.com>"
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build the Go app
RUN go build -o main .
# Expose port 80 to the outside world
EXPOSE 80
# Expose ENV Variabes
ENV NASA_API_KEY=""
ENV NASA_REDIS_URL="127.0.0.1:6379"
ENV NASA_PORT=":80"
# Command to run the executable
CMD ["./main"]