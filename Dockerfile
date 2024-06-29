FROM golang:1.22.4 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /build

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /build /build

EXPOSE 8080

CMD ["/build"]
