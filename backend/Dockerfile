FROM golang:latest as base

# Create another stage called "dev" that is based off of our "base" stage (so we have golang available to us)
FROM base as development

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

# print
COPY . .
COPY ./backend/.air.toml .
COPY ./backend/go.mod .
COPY ./backend/go.sum .

RUN go mod download


CMD ["air", "-c", ".air.toml"]
