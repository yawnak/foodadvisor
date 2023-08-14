FROM golang:1.21

RUN mkdir /app
WORKDIR /app

#Download dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY cmd ./cmd
COPY configs ./configs
COPY internal ./internal
COPY migrations ./migrations
COPY pkg ./pkg
COPY schema ./schema
COPY .env ./
COPY secrets.env ./


COPY . ./

RUN go build -o /build ./cmd/web

EXPOSE 8080

CMD [ "/build" ]