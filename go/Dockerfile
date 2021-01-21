FROM golang:1.14

WORKDIR /go/src/secretsanta
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["secretsanta"]