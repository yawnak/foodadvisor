FROM golang:1.21-alpine

RUN mkdir /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /build ./cmd/web

EXPOSE 8080

CMD [ "/build" ]