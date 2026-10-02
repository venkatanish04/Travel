# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o travelraft-server ./cmd/server

# Runtime stage
FROM scratch

WORKDIR /app

COPY --from=builder /app/travelraft-server .

EXPOSE 50051
EXPOSE 50052
EXPOSE 50053

CMD ["./travelraft-server"]
