FROM golang:1.24-alpine AS builder

# try with gco_enabled=0 when switch from sqlite to postgresql
RUN apk add --no-cache gcc g++ musl-dev
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -ldflags="-s -w" -o server

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/server ./server

ENTRYPOINT ["./server"]

