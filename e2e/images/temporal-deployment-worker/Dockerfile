# build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum files to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the worker binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker .

# final stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /app/worker /usr/local/bin/worker
ENTRYPOINT ["/usr/local/bin/worker"]
