FROM golang:1.25-alpine as build

# Create appuser.
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# Install Git and CA certificates
RUN apt-get update && apt-get install -y \
    ca-certificates \
    git

# Clone the hey repository
RUN git clone https://github.com/rakyll/hey.git /go/src/github.com/rakyll/hey

# Build
WORKDIR /go/src/github.com/rakyll/hey
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/hey hey.go

###############################################################################
# final stage
FROM alpine
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
USER 10001:10001

ARG APPLICATION="hey"
ARG DESCRIPTION="HTTP load generator, ApacheBench (ab) replacement, formerly known as rakyll/boom"
ARG PACKAGE="rakyll/hey"

LABEL org.opencontainers.image.ref.name="${PACKAGE}" \
    org.opencontainers.image.authors="Jaana Dogan <@rakyll>" \
    org.opencontainers.image.documentation="https://github.com/${PACKAGE}/README.md" \
    org.opencontainers.image.description="${DESCRIPTION}" \
    org.opencontainers.image.licenses="Apache 2.0" \
    org.opencontainers.image.source="https://github.com/${PACKAGE}"

COPY --from=build /go/bin/${APPLICATION} /hey
ENTRYPOINT ["/hey"]
