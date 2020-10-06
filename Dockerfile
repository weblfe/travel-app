FROM library/golang

# 启用 Go Modules 功能
ENV GO111MODULE on

# 配置 GOPROXY 环境变量
ENV GOPROXY https://goproxy.io

# Recompile the standard library without CGO
RUN CGO_ENABLED=0 go install -a std

# go path
ENV APP_DIR $GOPATH/src/github.com/weblfe/travel-app
RUN mkdir -p $APP_DIR


ADD . $APP_DIR

# Compile the binary and statically link
RUN cd $APP_DIR && CGO_ENABLED=0 go build -ldflags '-d -w -s'

EXPOSE 8080

# Set the entrypoint
ENTRYPOINT (cd $APP_DIR && ./travel-app)
