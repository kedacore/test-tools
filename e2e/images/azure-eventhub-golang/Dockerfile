FROM golang:1.18 as build-env

WORKDIR /go/src/app
ADD cmd/ /go/src/app
ADD go.* /go/src/app/

RUN go get -d -v ./...

RUN go build  -o /go/bin/app

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/app /
CMD ["/app"]