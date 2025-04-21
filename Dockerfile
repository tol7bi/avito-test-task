FROM golang:latest as builder

# Переходим в директорию приложения
WORKDIR /app
COPY . .



RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server .

RUN mkdir /root/logs

VOLUME ["/root/logs"]

RUN apk add --no-cache ca-certificates



ENV PORT=8080

EXPOSE 8080
EXPOSE 3000
EXPOSE 9000
ENTRYPOINT ["./server"]