FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
ARG TARGETARCH
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o /app/main .

# final stage
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
WORKDIR /
COPY --from=builder /app/main /app/main
# 65532 is numeric for nonroot
USER 65532:65532

ENTRYPOINT ["/app/main"]