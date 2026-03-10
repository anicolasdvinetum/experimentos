# build
FROM golang:1.25-alpine AS builder

WORKDIR /server
COPY go.mod go.sum ./
RUN go mod download

COPY . . 
RUN go build -o server

# runtime
FROM alpine:latest
WORKDIR /server
COPY --from=builder /server/server .
EXPOSE 8080

CMD ["./server"]