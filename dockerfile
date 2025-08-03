# --- Stage 1: Build the Go binary ---
  FROM golang:1.24.5

  # Set working directory
  WORKDIR /app
  
  # Copy all files including /v1 folder BEFORE go get
  COPY . .

  # Download deps
  RUN go mod tidy

  # Get dependencies
  RUN go get
  
  # Build the Go app
  RUN go build -o bin .
  
  EXPOSE 8080
  
  CMD ["./bin"]