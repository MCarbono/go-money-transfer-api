FROM golang:1.18

WORKDIR /app

COPY . /app

RUN go build -o main

EXPOSE 3000

CMD ["./main"]
