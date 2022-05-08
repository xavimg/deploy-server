FROM golang:alpine

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go get -u github.com/lib/pq

RUN go build -o main .

CMD ["./main"]