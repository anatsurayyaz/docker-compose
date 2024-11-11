FROM golang:1.22.3

# Set destination for COPY
WORKDIR /myapp

# Copy go.mod and execute to download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

COPY AUTHORS.md ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /my-go-app

COPY LINKS.md ./

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
EXPOSE 80

# Run
CMD [ "/my-go-app" ]

# Example how to ...
### docker build -t my-go-dockerfile .
### docker run -d -p 8085:8080 --name my-go-container my-go-dockerfile