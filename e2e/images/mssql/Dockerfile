FROM golang:1.24 as build-env

WORKDIR /go/src/app
COPY go.mod .
COPY cmd/ /go/src/app

RUN go mod tidy
RUN go build -o /go/bin/app

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/app /
CMD ["/app"]