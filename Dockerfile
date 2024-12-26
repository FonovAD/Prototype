FROM golang:1.23-bookworm
ENV CGO_ENABLED=1


RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      libsqlite3-0
WORKDIR /app
COPY go.* /app
RUN go mod download

COPY . /app

RUN GOOS=linux go build -o /app/main -v cmd/main.go

FROM golang:1.23-bookworm
ENV CGO_ENABLED=1

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      libsqlite3-0
COPY --from=0 /app /bin/app
WORKDIR /bin/app

CMD [ "./main" ]