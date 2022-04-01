FROM golang:1.17


WORKDIR /app
COPY go.mod .
COPY server.go server.go

RUN go mod tidy
RUN go build -o server .

EXPOSE 8080

CMD ["/app/server"]