# Define the builder stage
FROM golang:1.23.2-alpine3.20 AS builder
WORKDIR /app

# Copy and install go dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code and build the binary
COPY . .
RUN CGO_ENABLED=0 go build -o ./dist/ufogateway ./cmd/ufogateway/main.go

# Copy the binary to a new clean image
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/dist/ufogateway /app/ufogateway
CMD ["/app/ufogateway", "serve"]