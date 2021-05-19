FROM golang
MAINTAINER lixiangyun linimbus@126.com

RUN apt update && apt install gcc git -y

WORKDIR /gopath/

ENV GOPATH=/gopath/
ENV GOOS=linux
ENV BUILD_HOME=/gopath/src/github.com/easymesh/autoproxy/console

WORKDIR /gopath/src/github.com/easymesh

RUN git clone https://github.com/easymesh/autoproxy.git

RUN ls -ial

WORKDIR /gopath/src/github.com/easymesh/autoproxy/console
RUN go build -ldflags="-w -s" -o proxyweb

FROM ubuntu:20.04
MAINTAINER lixiangyun linimbus@126.com

RUN apt update && apt install libc-dev curl net-tools -y
ENV BUILD_HOME=/gopath/src/github.com/easymesh/autoproxy/console

WORKDIR /opt/
COPY --from=0 $BUILD_HOME/config.json ./config.json
COPY --from=0 $BUILD_HOME/release.db  ./release.db
COPY --from=0 $BUILD_HOME/proxyweb    ./proxyweb

RUN chmod +x proxyweb
EXPOSE 8000

ENTRYPOINT ["./proxyweb"]
