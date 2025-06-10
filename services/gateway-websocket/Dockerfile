FROM golang:1.24.3
WORKDIR /app
COPY . .
RUN go build -o /gateway

CMD ["/gateway"]
