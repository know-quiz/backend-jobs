# Start from the official Go image to create a build artifact
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

ARG GCP_CREDS_JSON_BASE64
ARG OPENAI_API_KEY
# Set the base64-encoded service account JSON as an environment variable
ENV GCP_CREDS_JSON_BASE64=$GCP_CREDS_JSON_BASE64
ENV OPENAI_API_KEY = $OPENAI_API_KEY

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]