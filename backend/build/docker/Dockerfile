FROM golang:1.23.0-alpine3.21

RUN apk --no-cache add curl

WORKDIR /go/arc/web

COPY . .

RUN go mod download

RUN go build -C ./cmd/server/

CMD ["./cmd/server/server"]