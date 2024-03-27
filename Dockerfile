FROM golang:alpine3.19
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o server ./cmd/web
CMD ["./server"]

