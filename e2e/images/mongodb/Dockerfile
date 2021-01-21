FROM golang:1.14 as build-env

WORKDIR /go/src/app
ADD cmd/ /go/src/app

RUN go get -d -v ./...

RUN go build -o /go/bin/app

FROM alpine:3.11.6
COPY --from=build-env /go/bin/app /
ENTRYPOINT ["/app"]