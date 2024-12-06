FROM golang:1.23.3 AS builder

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api .

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]
