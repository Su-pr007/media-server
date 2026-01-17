# Build stage
FROM golang:alpine AS builder
WORKDIR /src
COPY . .
RUN go build -o app

# Runtime stage
FROM alpine
WORKDIR /app
COPY --from=builder /src/app /app/
COPY public/index.html /app/public/index.html

EXPOSE 8080
ENTRYPOINT ["/app/app"]
