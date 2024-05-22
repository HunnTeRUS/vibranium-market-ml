FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /market-vibranium cmd/market-vibranium/main.go

EXPOSE 8080

CMD [ "/market-vibranium" ]
