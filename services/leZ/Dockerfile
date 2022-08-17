FROM golang:1.19.0-alpine as builder

WORKDIR /web-back-go/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY services/leZ/* ./


# pkg
COPY pkg ./pkg/

RUN go build -o /leZ


FROM alpine

COPY --from=builder /leZ /leZ
COPY services/leZ/.env /.env

RUN chmod +x /leZ

EXPOSE 9234

CMD ["/leZ"]