FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o w2w .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/w2w .

EXPOSE 8080

CMD ["./w2w"]


