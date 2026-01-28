FROM --platform=${BUILDPLATFORM} golang:1.25.6-alpine3.23 AS builder
WORKDIR $GOPATH/src/go.k6.io/k6
RUN apk --no-cache add git

ARG TARGETOS
ARG TARGETARCH
ARG K6_VERSION
RUN go install go.k6.io/xk6/cmd/xk6@latest
RUN GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" xk6 build \
        --k6-version "${K6_VERSION}" \
        --output /tmp/k6 \
        --with github.com/grafana/xk6-kubernetes@latest \
        --with github.com/grafana/xk6-disruptor@latest \
        --with github.com/JorTurFer/xk6-input-prometheus@latest

# Runtime 
FROM alpine:3.23 AS release

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 12345 -g 12345 k6
COPY --from=builder /tmp/k6 /usr/bin/k6

USER k6
WORKDIR /home/k6

ENTRYPOINT ["k6"]
