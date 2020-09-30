FROM jrottenberg/ffmpeg:4.1-alpine

WORKDIR /data/www/app

ADD ./app/api .
ADD ./conf  ./conf

RUN apk update --no-cache && apk add --no-cache ca-certificates && apk add --no-cache tzdata
ENV TZ Asia/Shanghai

VOLUME /data/www/app/static

EXPOSE 8080

ENTRYPOINT [ "/data/www/app/api" ]
