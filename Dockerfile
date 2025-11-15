FROM golang:1.23-alpine AS builder

RUN echo "https://mirror.arvancloud.ir/alpine/v3.22/main" > /etc/apk/repositories && \
    echo "https://mirror.arvancloud.ir/alpine/v3.22/community" >> /etc/apk/repositories

RUN apk update && apk add --no-cache git

WORKDIR /app

RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api

FROM alpine:latest

RUN echo "https://mirror.arvancloud.ir/alpine/v3.22/main" > /etc/apk/repositories && \
    echo "https://mirror.arvancloud.ir/alpine/v3.22/community" >> /etc/apk/repositories

RUN apk update && apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]