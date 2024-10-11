FROM golang:1.22
WORKDIR /cmd

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o cmd .

ENTRYPOINT ["./cmd"]
