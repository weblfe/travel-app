FROM golang:1.15.0-alpine3.12

WORKDIR /data/www/app

COPY .  ./

RUN  go env -w GOPROXY=https://goproxy.io \
    && export GO111MODULE=on && export GOPROXY=https://goproxy.io \
    && go env -w GOPRIVATE=*.gitlab.com,*.gitee.com \
    && go build  -o api  \
    && rm -rf common controllers libs middlewares models plugins repositories routers services  \
      tests transfroms transports .env data go.* .gitigrnore deploy.sh lastupdate.tmp \
    && rm -rf ${GOPATH}/src  && rm -rf ${GOPATH}/pkg && rm -rf ${GOPATH}/mod

EXPOSE 8080
EXPOSE 8088

VOLUME /data/www/app/.env
VOLUME /data/www/app/static

ENTRYPOINT [ "/data/www/app/api" ]