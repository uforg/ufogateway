FROM golang:1.23.2
WORKDIR /app

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the files and build the app
COPY . .
RUN go build -o ./dist/ufogateway ./cmd/ufogateway/main.go

# Run the app
EXPOSE 8080
CMD ["./dist/ufogateway"]
