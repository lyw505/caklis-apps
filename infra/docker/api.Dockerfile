FROM golang:1.25-alpine

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Expose port and run Air
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]
