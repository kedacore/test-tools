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

FROM alpine:3.18
RUN apk add --no-cache ca-certificates && \
    adduser -D -u 65532 -g 65532 k6
COPY --from=builder /tmp/k6 /usr/bin/k6
USER 65532
WORKDIR /home/k6

ENTRYPOINT ["k6"]
