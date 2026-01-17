# Build stage
FROM golang:alpine AS builder
WORKDIR /src
COPY . .
RUN go build -o app

# Runtime stage
FROM alpine
COPY --from=builder /src/app /app/
COPY public /app/public

EXPOSE 8080
ENTRYPOINT ["/app/app"]
