FROM ubuntu:20.04
MAINTAINER lixiangyun linimbus@126.com

RUN apt update && apt install golang gcc -y

WORKDIR /gopath/
ENV GOPATH=/gopath/
ENV GOOS=linux
ENV BUILD_HOME=/gopath/src/github.com/easymesh/autoproxy/console

COPY . /gopath/src/github.com/easymesh/autoproxy/console

WORKDIR /gopath/src/github.com/easymesh/autoproxy/console
RUN go build -ldflags="-w -s" -o proxyweb

FROM ubuntu:20.04
MAINTAINER lixiangyun linimbus@126.com

RUN apt update && apt install libc-dev curl net-tools

WORKDIR /opt/
COPY --from=0 $BUILD_HOME/config.json ./config.json
COPY --from=0 $BUILD_HOME/release.db  ./release.db
COPY --from=0 $BUILD_HOME/proxyweb    ./proxyweb

RUN chmod +x proxyweb
EXPOSE 8000

ENTRYPOINT ["./proxyweb"]
