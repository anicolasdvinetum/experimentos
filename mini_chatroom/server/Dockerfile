FROM golang:1.21-alpine AS builder

WORKDIR /mini_chatroom
COPY go.mod ./
RUN go mod download

COPY . . 
RUN go build -o server


FROM alpine:latest
WORKDIR /mini_chatroom
COPY --from=builder /mini_chatroom/server .
EXPOSE 9000

CMD ["./server"]
