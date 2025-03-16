# Build stage
FROM golang:1.22.5 AS builder
WORKDIR /app
COPY . .
RUN go build -o location-api ./cmd

# Run stage
FROM ubuntu:latest
LABEL authors="can.zanat"
WORKDIR /app/
COPY --from=builder /app/location-api ./
COPY ./.config ./.config
RUN chmod +x location-api
CMD ["./location-api"]