FROM node:20-alpine AS svelte-builder
WORKDIR /app

COPY client_web/package*.json .
RUN npm ci

COPY client_web/ .
RUN npm run build

FROM golang:1.25-alpine AS go-builder
WORKDIR /app

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ .
RUN go build -o server -ldflags="-s -w" main.go

FROM alpine:3.20
RUN apk add --no-cache ffmpeg

WORKDIR /app
COPY --from=go-builder /app/server ./server/
COPY --from=svelte-builder /app/build ./client_web/build

ENV GIN_MODE=release

WORKDIR /app/server
ENTRYPOINT ["./server"]

