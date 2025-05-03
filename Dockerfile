FROM golang:alpine AS builder

WORKDIR /calc

RUN apk add --no-cache  \
    git \
    gcc \ 
    musl-dev \
    sqlite-dev

COPY . .

RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 go build -x -o calc ./cmd

FROM alpine:latest

WORKDIR /

COPY --from=builder /calc/calc ./calc

EXPOSE 8080

CMD [ "./calc" ]