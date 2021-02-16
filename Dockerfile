FROM golang:1.15-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Install the package
RUN go install -v ./...

# Run the executable
# CMD ["open-func"]