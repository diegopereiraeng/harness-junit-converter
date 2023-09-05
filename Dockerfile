FROM golang:1.21-alpine as build

# Install git for fetching dependencies
RUN apk add --no-cache --update git

# Set working directory inside the container
WORKDIR /app

# Copy go files and go.mod to working directory
COPY *.go ./
COPY go.mod go.sum ./

# Fetch dependencies
RUN go mod download

# Build the Go application
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o harness-junit-converter

# Final stage
FROM alpine:3.14

# Copy the binary from the build stage to the final stage
COPY --from=build /app/harness-junit-converter /bin/

# Set working directory
WORKDIR /bin

# Set the binary as the entry point of the container
ENTRYPOINT ["/bin/harness-junit-converter"]
