FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /build
COPY . .

RUN go mod download
RUN go build -o todo_api .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /build/todo_api .

EXPOSE 8080

CMD ["./todo_api"]