FROM golang:1.20-alpine AS builder
RUN apk add --no-cache git ffmpeg build-base libmediainfo-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=1
RUN go build -o server .

FROM alpine:3.17
RUN apk add --no-cache ffmpeg ca-certificates mediainfo
COPY --from=builder /app/server /server
ENTRYPOINT ["/server"]