FROM golang:1.16-alpine3.15  AS builder
LABEL stage=builder

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build infura-test

FROM alpine:latest AS production
COPY --from=builder /app .

ENV API_BASE_URL=https://mainnet.infura.io/v3/ \
    INFURA_PROJECT_ID=56342c82e68e4d11bc998e0a54d16d84 \
    PORT=8000

CMD ["./infura-test"]

# build me with `docker build -t infura-test ."
# run me with `docker run -p 8000:8000 infura-test`
