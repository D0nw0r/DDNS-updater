# Use a base image for your Go app
FROM golang:1.23.2-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Copy the .env file into the container (make sure the path is correct)
COPY readenv/.env /app/readenv/.env

# Install dependencies and build your Go app (adjust as needed)
RUN go mod tidy
RUN go build -o myapp .

# Set the entry point for the app (replace with your app's executable)
CMD ["./myapp"]
