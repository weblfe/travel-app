FROM golang:1.16-alpine AS builder

# 启用 Go Modules 功能
ENV GO111MODULE on

# 配置 GOPROXY 环境变量
ENV GOPROXY "https://goproxy.cn,direct"

# set GOPATH
ENV GOPATH /go

# Recompile the standard library without CGO
RUN CGO_ENABLED=0 go install -a std

# go path
ENV APP_DIR $GOPATH/src/github.com/weblfe/travel-app
RUN mkdir -p $APP_DIR

ADD . $APP_DIR

# Compile the binary and statically link
RUN cd $APP_DIR  && export version=$(/bin/date "+%Y-%m-%d %H:%M:%s") && CGO_ENABLED=0 GOOS=linux go build -ldflags="-d -w -s" -ldflags="-X 'main.BuildTime=${version}'" -o  travel-app main.go
RUN cd $APP_DIR/cmder  && export version=$(/bin/date "+%Y-%m-%d %H:%M:%s") && CGO_ENABLED=0 GOOS=linux go build -ldflags="-d -w -s" -ldflags="-X 'main.BuildTime=${version}'" -o  travel-tools main.go

FROM jrottenberg/ffmpeg:4.1-alpine

WORKDIR /data/www/app

RUN apk update --no-cache && apk add --no-cache ca-certificates && apk add --no-cache tzdata
ENV TZ Asia/Shanghai

ADD ./conf                         /data/www/app/conf
ADD ./entrypoint-api.sh            /data/www/app/entrypoint.sh
COPY --from=builder /go/src/github.com/weblfe/travel-app/travel-app    /data/www/app/api-server
COPY --from=builder /go/src/github.com/weblfe/travel-app/cmder/travel-tools    /data/www/app/api-cli

VOLUME /data/www/app/static
VOLUME /data/www/app/views

EXPOSE 8080

RUN chmod +x /data/www/app/api-server && chmod +x /data/www/app/api-cli

# Set the entrypoint
ENTRYPOINT (cd /data/www/app && ./api-server)