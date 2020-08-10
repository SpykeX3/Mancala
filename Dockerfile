FROM golang:1.14

ENV GO111MODULE=on

WORKDIR /go/src/app
COPY . .
COPY .env /

RUN go mod download
RUN go build -o /server /go/src/app/Server/

EXPOSE 1337

CMD ["/server"]