FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o article-tag-extractor ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/article-tag-extractor .

CMD ["./article-tag-extractor"]