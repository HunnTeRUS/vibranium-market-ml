FROM golang:1.20

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/market-vibranium cmd/market-vibranium/main.go

EXPOSE 8080

RUN mkdir -p /app/snapshots

ENTRYPOINT ["/app/market-vibranium"]