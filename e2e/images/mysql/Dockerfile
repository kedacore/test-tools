# build stage
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
ARG TARGETARCH
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o /app/main ./cmd/main.go

# final stage
FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /app/main /app/main
ENTRYPOINT ["/app/main"]