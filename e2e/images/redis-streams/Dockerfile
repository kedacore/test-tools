FROM golang:1.25-alpine as builder
WORKDIR /cmd
COPY ./cmd .
COPY go.mod go.sum ./
RUN go build -tags netgo -ldflags '-w -s' -o main .
CMD [ "./main" ]

FROM scratch
COPY --from=builder /cmd/main /main
CMD ["/main"]