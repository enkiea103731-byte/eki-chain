FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN go build -o eki

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /app/eki .
CMD ["./eki"]
