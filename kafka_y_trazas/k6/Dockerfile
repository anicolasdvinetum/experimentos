FROM golang:1.21 AS builder

# instalar xk6 v0.12.0 compatible con k6 v0.49.0
RUN go install go.k6.io/xk6/cmd/xk6@v0.12.0

WORKDIR /src

RUN xk6 build v0.49.0 \
    --with github.com/mostafa/xk6-kafka@v0.21.0 \
    --output ./k6

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/k6 /usr/bin/k6

ENTRYPOINT ["k6"]