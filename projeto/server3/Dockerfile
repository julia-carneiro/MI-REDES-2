FROM golang:1.20

WORKDIR /app

COPY ../go.mod ../go.sum ./
COPY . .

RUN go build -o server .

CMD ["./server"]
