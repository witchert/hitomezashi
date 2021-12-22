FROM golang:1.16

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /docker-hitomezashi

RUN chmod +x /docker-hitomezashi

EXPOSE 1111

CMD [ "/docker-hitomezashi" ]