FROM golang:1.22-alpine AS builder 

WORKDIR /myapp

COPY go.mod go.sum ./ 

RUN go mod download 

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/myapp ./cmd/main.go

FROM alpine:latest AS runner  

RUN apk add --no-cache  ca-certificates 

WORKDIR /root/

COPY --from=builder /usr/local/bin/myapp .

COPY configs/config.yaml ./configs/

EXPOSE 8080

 CMD ["./myapp"]