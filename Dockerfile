# docker build -t collapp:v1 .
# docker run -d -p 3000:3000 --name collapp collapp:v1

FROM golang:1.18.2

WORKDIR /app

COPY . .

RUN go build -o collapp

EXPOSE 8080

CMD ./collapp 