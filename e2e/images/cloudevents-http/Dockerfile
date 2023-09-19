FROM golang:1.19 as build-env


WORKDIR /app
COPY go.mod .
COPY server.go server.go

RUN go mod tidy
RUN go build -o /go/bin/server .

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/server /
EXPOSE 8899

ENTRYPOINT ["/server"]