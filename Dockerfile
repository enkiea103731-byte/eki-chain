FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go mod tidy && go build -o eki-chain .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/eki-chain .
EXPOSE 8080
CMD ["./eki-chain"]
