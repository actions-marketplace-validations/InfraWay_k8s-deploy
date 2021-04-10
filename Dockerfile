FROM golang:1.15.5-alpine3.12
WORKDIR /usr/src/app
COPY . /usr/src/app/
RUN go mod download && go build .
ENTRYPOINT ["/usr/src/app/deploy"]