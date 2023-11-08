FROM golang:1.18


WORKDIR /crud
COPY go.* ./
RUN go mod download


CMD ["go", "run", "cmd/bin/server/main.go"]

