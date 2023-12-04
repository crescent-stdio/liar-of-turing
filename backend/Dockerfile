# FROM golang:1.21.4-bullseye as base

# FROM base as development

# WORKDIR /app

# RUN go install github.com/cosmtrek/air@latest

# COPY ./backend/go.mod .
# COPY ./backend/go.sum .
# RUN go mod download

# CMD ["air", "-c", ".air.toml"]
FROM golang:1.21.4-bullseye as base

FROM base as development

WORKDIR /usr/src/app

COPY ./backend/go.mod .
COPY ./backend/go.sum .

RUN go mod tidy

COPY ./backend/ .

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -a -ldflags="-s -w" -o bin/main cmd/server/main.go

FROM debian:bullseye

COPY --from=base /usr/src/app/bin/main ./main

ENV PORT 8080
ENV HOSTNAME "127.0.0.1"
EXPOSE 8080

RUN chmod +x ./main

ENTRYPOINT ["./main"]
