# Multi-stage build for minimal backend container size
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

FROM alpine:3.20

WORKDIR /app

# Copy compiled binary and resources
COPY --from=builder /app/server .
COPY --from=builder /app/resources ./resources

EXPOSE 7860 8081 8082

CMD ["./server"]
