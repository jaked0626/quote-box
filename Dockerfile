# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -v -o server ./cmd/web

# Run stage 
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
ENV PORT 4000
EXPOSE 4000
ENTRYPOINT [ "/app/server" ]

