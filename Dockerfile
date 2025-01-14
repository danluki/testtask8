FROM golang:1.23.4 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./main .

FROM alpine:3.18

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main /bin/main

EXPOSE 8080

CMD ["/bin/main"]