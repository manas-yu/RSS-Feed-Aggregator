FROM golang:1.21.4-alpine

# Set the working directory
WORKDIR /app

# Install git (required to install Goose)
RUN apk add --no-cache git

# Install Goose
RUN go install github.com/pressly/goose/cmd/goose@latest

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Run migrations and then start the application
CMD goose -dir=/app/sql/schema postgres "postgres://postgres:abj1195@db:5432/rssagg?sslmode=disable" up && ./main