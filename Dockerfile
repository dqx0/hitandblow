# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o hitandblow ./cmd/main.go

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/hitandblow .
EXPOSE 8080
CMD ["./hitandblow"]