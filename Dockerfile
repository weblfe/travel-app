FROM golang:1.15-alpine3.12 AS builder

WORKDIR /data/www/app

COPY . ./

RUN  go env -w GOPROXY=https://goproxy.io \
    && export GO111MODULE=on && export GOPROXY=https://goproxy.io \
    && go env -w GOPRIVATE=*.gitlab.com,*.gitee.com \
    && go build -ldflags="-s -w" -o /data/www/app/api

FROM alpine

WORKDIR /data/www/app

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates && apk add --no-cache tzdata
ENV TZ Asia/Shanghai

COPY --from=builder /data/www/app/api    /data/www/app/api
COPY --from=builder /data/www/app/conf    /data/www/app/conf

EXPOSE 8080
EXPOSE 8088

VOLUME /data/www/app/.env
VOLUME /data/www/app/static

ENTRYPOINT [ "/data/www/app/api" ]
