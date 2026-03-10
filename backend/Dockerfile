# Build stage
FROM golang:1.25-alpine AS builder

RUN apk --no-cache add gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o eventdoctor ./cmd/api/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

COPY --from=builder /app/eventdoctor .

EXPOSE 8087

CMD ["./eventdoctor"]
