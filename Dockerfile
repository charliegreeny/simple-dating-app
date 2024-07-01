FROM golang:1.22.4

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . . 

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o exe ./cmd/main.go

EXPOSE 8080

# Run
CMD ["/app/exe"]