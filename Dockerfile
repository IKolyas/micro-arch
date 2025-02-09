FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

FROM scratch

WORKDIR /build

COPY --from=builder /app/main /build/main

ENTRYPOINT ["/build/main"]