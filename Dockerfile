# Start from the latest golang base image
FROM golang:latest as builder

LABEL maintainer="rtkennelly1@gmail.com"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./


# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o livetest ./cmd/livetest

######## Start a new stage from scratch #######
FROM scratch

WORKDIR /root/

# Copy the config and pre-built binary file from the previous stage
COPY --from=builder /app/config.json .
COPY --from=builder /app/livetest .

# Command to run the executable
CMD ["./livetest"]
