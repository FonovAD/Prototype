FROM golang:1.23.3-alpine3.20

WORKDIR /app
COPY go.* /app
RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main -v cmd/main.go

FROM alpine:3.20

COPY --from=0 /app /bin/app
WORKDIR /bin/app

CMD [ "./main" ]
