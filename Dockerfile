FROM golang:1.22
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o eki-chain .
EXPOSE 8080
CMD ["./eki-chain"]
