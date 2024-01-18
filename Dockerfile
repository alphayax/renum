# Start from the latest golang base image
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the builder container
WORKDIR /app

# Copy go mod and sum files
COPY go.* ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the builder container
COPY . .

# Build the Go app as a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o renum .

# Start a new stage from scratch
FROM scratch

# Copy the binary from builder
COPY --from=builder /app/renum .

# Command to run the executable
ENTRYPOINT ["./renum"]
