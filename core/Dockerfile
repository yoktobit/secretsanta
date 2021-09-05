FROM golang:latest AS build

WORKDIR /app
COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o secretsanta .

FROM alpine:latest
COPY --from=build /app/secretsanta /app/secretsanta

CMD ["/app/secretsanta"]