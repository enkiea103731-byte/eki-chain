FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY *.go ./
COPY *.json ./
RUN go build -o eki .
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/eki ./
COPY --from=build /app/*.json ./
COPY site ./site
COPY index.html ./
EXPOSE 10000
CMD ["./eki"]
