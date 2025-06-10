FROM golang:1.24.3
WORKDIR /app
COPY . .
RUN go build -o /message-service .

CMD ["/message-service"]