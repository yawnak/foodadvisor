FROM golang:1.21-alpine

RUN mkdir /app
WORKDIR /app

COPY cmd ./cmd
COPY configs ./configs
COPY internal ./internal
COPY migrations ./migrations
COPY pkg ./pkg
COPY schema ./schema
COPY .env ./
COPY secrets.env ./
COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN go build -o /build ./cmd/web

EXPOSE 8080

CMD [ "/build" ]