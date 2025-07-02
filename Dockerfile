FROM golang:1.20-alpine AS builder
RUN apk add --no-cache git ffmpeg
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM alpine:3.17
RUN apk add --no-cache ffmpeg ca-certificates
COPY --from=builder /app/server /server
ENTRYPOINT ["/server"]