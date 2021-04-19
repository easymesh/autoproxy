export CGO_ENABLED=1
export GARCH=$(go env GOARCH)

go build -ldflags="-w -s" -o proxyweb
tar -zcf proxyweb_linux_$GARCH.tar.gz.tar.gz proxyweb release.db config.json
