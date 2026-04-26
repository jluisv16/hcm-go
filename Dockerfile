FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/hcm-api ./cmd/api

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app

USER app
WORKDIR /app

COPY --from=builder /bin/hcm-api /app/hcm-api

EXPOSE 8080

ENV HTTP_PORT=8080
ENV APP_ENV=production

ENTRYPOINT ["/app/hcm-api"]
