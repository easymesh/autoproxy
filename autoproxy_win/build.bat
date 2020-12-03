rsrc -manifest exe.manifest -ico static/main.ico
rice embed-go
go build -ldflags="-H windowsgui -w -s" -o autoproxy_desktop.exe
zip autoproxy_desktop.zip autoproxy_desktop.exe
pause