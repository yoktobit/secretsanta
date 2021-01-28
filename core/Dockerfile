FROM golang:alpine

WORKDIR /go/src/secretsanta
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["secretsanta"]