FROM golang:1.17 as build-env

WORKDIR /app
COPY . .
RUN go mod tidy

RUN go build -o /go/bin/server .
RUN go build -o /go/bin/client ./cmd

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/server /
COPY --from=build-env /go/bin/client /
EXPOSE 8080

CMD ["/server"]