# STILL WORK IN PROGRESS
FROM node:lts-buster AS frontend_builder

COPY ./frontend /frontend

WORKDIR /frontend

RUN set -ex \
    && yarn --network-timeout 1800000 \
    && yarn build


FROM golang:alpine AS backend_builder

ENV GO111MODULE on

COPY . /Drive
COPY --from=frontend_builder /frontend/build/ /Drive/frontend/build/

WORKDIR /Drive

RUN set -ex \
    && apk upgrade \
    && apk add gcc libc-dev git \
    && go get github.com/rakyll/statik \
    && go generate \
    && go install \
    && go build -o "Drive"

FROM alpine AS final

ARG TZ="Asia/Shanghai"

ENV TZ ${TZ}

COPY --from=backend_builder /Drive/Drive /Drive/Drive

RUN apk upgrade \
    && apk add bash tzdata \
    && ln -s /Drive/Drive /usr/bin/Drive \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && rm -rf /var/cache/apk/*

EXPOSE 8421/tcp

ENTRYPOINT ["Drive"]