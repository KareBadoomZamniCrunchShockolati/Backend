FROM golang:1.23.0-alpine AS builder

RUN apk add --no-cache gcc musl-dev libwebp-dev jpeg-dev

WORKDIR /app

RUN go env -w GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/bidlancer

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata libwebp-dev jpeg-dev

WORKDIR /app

RUN mkdir -p ./internal/jwt ./SSL ./internal/application/email_templates

COPY --from=builder /app/main .



EXPOSE 8080

CMD ["./main"]