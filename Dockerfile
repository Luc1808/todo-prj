# 1) Builder stage: compile the Go binary
FROM golang:1.24-alpine AS builder

# (optional) install git if you need private modules
# RUN apk add --no-cache git

WORKDIR /app
# Copy module files and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source, then build
COPY . .
# Adjust -o name to whatever you like
RUN go build -o api ./app

# 2) Final stage: a minimal runtime image
FROM alpine:3.18

# (optional) add CA certs if your API makes TLS calls
RUN apk add --no-cache ca-certificates

WORKDIR /app
# Copy the compiled binary
COPY .env .
COPY --from=builder /app/api .

# Expose the port your app listens on
EXPOSE 8080

# Launch!
CMD ["./api"]
