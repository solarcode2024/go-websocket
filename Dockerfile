# Use the official GoLang base image from Docker Hub
FROM golang:1.21.0

# Set the working directory inside the container
WORKDIR /sportnex-websocket

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY . .

# Build the GoLang application inside the container
# RUN CGO_ENABLED=0 GOOS=linux go build -o bin/lpi ./cmd

# Set the default command to run the GoLang application
CMD ["./bin/sportnex-websocket"]
