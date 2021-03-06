FROM golang:1-alpine as builder
RUN set -ex \
    && apk add --no-cache --virtual .build-deps \
    gcc libc-dev git
COPY ./src /src
WORKDIR /src
RUN go mod tidy && go mod download && go build -o /bin/docker-plugin-volume

FROM alpine:3
COPY --from=builder /bin/docker-plugin-volume /bin/docker-plugin-volume
CMD ["/bin/docker-plugin-volume"]