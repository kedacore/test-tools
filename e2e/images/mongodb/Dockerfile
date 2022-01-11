FROM golang:1.17 as build-env

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

ADD cmd cmd

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app cmd/main.go

FROM gcr.io/distroless/static:nonroot
COPY --from=build-env /go/bin/app /
USER nonroot:nonroot

ENTRYPOINT ["/app"]