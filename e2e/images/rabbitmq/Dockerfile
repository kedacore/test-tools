# build stage
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
ARG TARGETARCH
WORKDIR /app

# Copy go.mod and go.sum files to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binaries for send and receive
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o /app/send ./cmd/send
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o /app/receive ./cmd/receive

# final stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /app/send /usr/local/bin/
COPY --from=builder /app/receive /usr/local/bin/

# By default, the image won't do anything.
# The user will specify which binary to run.
CMD ["receive"]