FROM golang:1.19.0-alpine as builder

WORKDIR /web-back-go/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY services/tinyURL/* ./


# pkg
COPY pkg ./pkg/

RUN go build -o /tinyURL


FROM alpine

COPY --from=builder /tinyURL /tinyURL
COPY services/tinyURL/.env /.env

RUN chmod +x /tinyURL

EXPOSE 9234

CMD ["/tinyURL"]