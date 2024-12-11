FROM golang:1.20 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o service-discovery .

FROM alpine:3.18

WORKDIR /root/
COPY --from=builder /app/service-discovery .
EXPOSE 8080
CMD ["./service-discovery"]
