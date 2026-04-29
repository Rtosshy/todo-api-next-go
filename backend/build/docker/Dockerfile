FROM golang:1.26.2-alpine3.23 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.23

RUN apk --no-cache add ca-certificates curl

WORKDIR /app
COPY --from=builder /out/server /app/server

EXPOSE 8080

CMD ["/app/server"]
