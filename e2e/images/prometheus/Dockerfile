FROM golang:1.19 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go

RUN go build -o prometheusapp main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
WORKDIR /
COPY --from=builder /workspace/prometheusapp .
# 65532 is numeric for nonroot
USER 65532:65532

ENTRYPOINT ["/prometheusapp"]