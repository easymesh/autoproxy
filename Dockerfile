FROM golang

ENV GOPATH=/go/
ENV CGO_ENABLED=0
ENV BUILD_HOME=/go/src/github.com/easymesh/autoproxy

WORKDIR $BUILD_HOME

COPY . .

RUN go build -ldflags="-w -s" -o autoproxy

FROM scratch

ENV BUILD_HOME=/go/src/github.com/easymesh/autoproxy
WORKDIR /bin

COPY --from=0 $BUILD_HOME/domain.json /bin
COPY --from=0 $BUILD_HOME/autoproxy   /bin

EXPOSE 8000-9000

ENTRYPOINT ["/bin/autoproxy"]