FROM golang:1.21.4-alpine as base

FROM base as development

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV PORT 8080
ENV HOSTNAME "0.0.0.0"

CMD ["air", "-c", ".air.toml"]

###############################

# Use golang:1.21.4-alpine as the base image
# FROM golang:1.21.4-alpine as base

# # Start the development stage based on the base image
# FROM base as development

# # Set the working directory inside the container
# WORKDIR /app

# # Copy the go.mod and go.sum files first to leverage Docker cache
# COPY go.mod .
# COPY go.sum .

# # Tidy up the go modules
# RUN go mod tidy

# COPY . .
# # Copy the rest of the backend directory

# # Set environment variables for Go
# ENV GO111MODULE=on \
#     CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# # Build the Go application
# RUN go build -a -ldflags="-s -w" -o bin/main cmd/server/main.go

# # Start a new stage from alpine for a smaller final image
# FROM alpine

# # Install necessary dependencies for the runtime
# RUN apk add --no-cache ca-certificates
# RUN apk update
# RUN apk add git
# RUN apk add upx

# # Copy the binary from the development stage to the current stage
# COPY --from=development /app/.env ./.env
# COPY --from=development /app/data ./data
# COPY --from=development /app/bin/main ./main

# # Set environment variables
# ENV PORT 8080
# ENV HOSTNAME "0.0.0.0"

# # Expose port 8080
# EXPOSE 8080

# # Make the binary executable
# RUN chmod +x ./main

# # Define the entrypoint
# ENTRYPOINT ["./main"]
