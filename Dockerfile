FROM golang:1.18-alpine as build

RUN apk update && apk add git

WORKDIR /app

ADD go.mod go.sum ./

RUN go mod download

ADD ./ ./
ENV CGO_ENABLED=0
RUN go build

FROM alpine as runtime

COPY --from=build /app/akv-entrypoint /usr/bin

ENTRYPOINT ["/usr/bin/akv-entrypoint"]