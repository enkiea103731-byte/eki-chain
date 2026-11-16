FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o eki main.go
EXPOSE 10000
CMD ["./eki"]
