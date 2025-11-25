FROM golang:1.25-alpine as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go

RUN go build -o boundserviceaccounttoken main.go

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
WORKDIR /
COPY --from=builder /workspace/boundserviceaccounttoken .
# 65532 is numeric for nonroot
USER 65532:65532
EXPOSE 8080

ENTRYPOINT ["/boundserviceaccounttoken"]
