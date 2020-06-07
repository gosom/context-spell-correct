FROM golang:latest as builder
LABEL maintainer="Giorgos Komninos <info@zendrom.com>"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root

COPY --from=builder /app/main .

CMD ["./main"] 
