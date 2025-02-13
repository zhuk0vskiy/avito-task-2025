FROM golang:1.23.2

RUN set -ex
USER root

WORKDIR /app

COPY . /app

RUN go build -o /build ./cmd \
    && go clean -cache -modcache

RUN 

ENV CONFIG_PATH=/app/config/config.yml
ENV JWT_KEY=OhGodILoveAvito

CMD ["/build"]