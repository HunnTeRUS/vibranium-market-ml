FROM golang:1.20

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

COPY ./cmd/sql/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

RUN go build -o /market-vibranium cmd/market-vibranium/main.go

EXPOSE 8080

VOLUME /snapshots

ENTRYPOINT ["/wait-for-it.sh", "db:3306", "--"]
CMD ["./market-vibranium"]