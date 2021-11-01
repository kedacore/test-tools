FROM golang

WORKDIR /cmd

COPY ./cmd .
COPY go.mod go.sum ./

RUN go build -o main .

CMD [ "./main" ]