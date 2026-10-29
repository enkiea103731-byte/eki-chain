FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy && go build -o eki-chain .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/eki-chain .
COPY --from=builder /app/*.json ./
EXPOSE 8080
CMD ["./eki-chain"]
