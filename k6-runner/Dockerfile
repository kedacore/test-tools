FROM --platform=${BUILDPLATFORM} golang:1.21-alpine3.18 as builder
WORKDIR $GOPATH/src/go.k6.io/k6
RUN apk --no-cache add git

ARG TARGETOS
ARG TARGETARCH
RUN go install go.k6.io/xk6/cmd/xk6@latest
RUN GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" xk6 build \
        --output /tmp/k6 \
        --with github.com/grafana/xk6-kubernetes@latest \
        --replace github.com/grafana/xk6-kubernetes=github.com/JorTurFer/xk6-kubernetes@latest\
        --with github.com/grafana/xk6-disruptor@latest \
        --with github.com/JorTurFer/xk6-input-prometheus@latest

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /tmp/k6 /usr/bin/k6
# 65532 is numeric for nonroot
USER 65532:65532

ENTRYPOINT ["k6"]
