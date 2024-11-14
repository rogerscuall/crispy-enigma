# From alpine lates load the binary in dist folder
FROM alpine:latest
COPY dist/crispy-enigma_linux_386/crispy-enigma /crispy-enigma
CMD sh

# Start from the official Go image.
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the downloaded repo code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 go build -o crispy-enigma .

# Use a smaller base image if possible (optional but recommended)
FROM alpine:latest

# Copy the compiled binary from the builder stage
COPY --from=0 /app/crispy-enigma /crispy-enigma

# Set entrypoint to run your Go binary
CMD ["sh"]
