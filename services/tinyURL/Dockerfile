FROM golang:1.19.0-alpine

WORKDIR /web-back-go/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY services/tinyURL/* ./


# pkg
COPY pkg ./pkg/

RUN go build -o /tinyURL

EXPOSE 9234

CMD ["/tinyURL"]