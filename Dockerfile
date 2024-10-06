FROM golang:1.22.7 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /build

FROM alpine:3.20.3

RUN apk --no-cache add ca-certificates

COPY --from=builder /build /build

EXPOSE 80
EXPOSE 443

CMD ["/build"]
