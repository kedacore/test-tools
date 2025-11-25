FROM golang:1.25 as build-env

# Create MQ installation
RUN mkdir -p /opt/mqm \
    && chmod a+rx /opt/mqm
# Location of the downloadable MQ client package \
ENV RDURL="https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqdev/redist" \
    RDTAR="IBM-MQC-Redist-LinuxX64.tar.gz" \
    VRMF=9.2.0.0
# Install the MQ client from the Redistributable package. This also contains the
# header files we need to compile against. Setup the subset of the package
# we are going to keep - the genmqpkg.sh script removes unneeded parts
ENV genmqpkg_incnls=1 \
    genmqpkg_incsdk=1 \
    genmqpkg_inctls=1
RUN cd /opt/mqm \
    && curl -LO "$RDURL/$VRMF-$RDTAR" \
    && tar -zxf ./*.tar.gz \
    && rm -f ./*.tar.gz \
    && bin/genmqpkg.sh -b /opt/mqm

WORKDIR /go/src/app
COPY go.mod .
COPY cmd/ /go/src/app
RUN go mod tidy
RUN go build -o /go/bin/app


# We use the `distroless/cc` because the `mq-golang` client depends on `libstdc++6` 
FROM gcr.io/distroless/cc
COPY --from=build-env /go/bin/app /
COPY --from=build-env /opt/mqm /opt/mqm
CMD ["/app"]