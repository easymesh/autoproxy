# Autoproxy

The project provides browser http proxy proxy service, supports https and http protocol proxy; can deploy multi-level proxy; supports TLS protocol encryption; mainly uses tools for accessing external websites with intranet hosts;

## Build
- Prepare linux amd64 and golang.
- `git clone https://github.com/easymesh/autoproxy.git`
- `cd autoproxy`
- `bash build.sh`

## Features
- Support multiple forwarding modes
- Support statistical Console
- Support secondary forwarding
- Support TLS transmission encryption
- Support dynamic routing
- Support custom forwarding domain name settings
- Support multiple platforms

## Deploy

### RemoteServer(Not-desktop)
- Prepare a virtual cloud host with public IP
- Specify a port for open tcp protocol access, such as 8080

### [Download Binary](https://github.com/easymesh/autoproxy/releases)

- Choose the right platform, Download the latest version；such as. `autoproxy_linux_amd64.tar.gz`
- Run `unzip autoproxy_linux_amd64.tar.gz` Unzip the compressed package.
- Run `nohup xxx &` The program will run in the background.

- Remote Deploy:
```
./autoproxy -local-address https://0.0.0.0:8080 -mode local -local-auth user:123456
```

- Local Deploy:
```
./autoproxy -local-auth "" -remote-address https://{remote-ip}:8080 -remote-auth user:123456
```

### Default configuration
- The default current path is the log storage path
- The default bound port 8080
- TLS transmission encryption is enabled by default
- Provide two default authentication accounts

```
Usage of ./autoproxy:
  -cert-file string
        tls cert file pem format, if not set, the program will automatically generate
  -domain string
        match domain list file(domain mode requires) (default "domain.json")
  -help
        usage help
  -key-file string
        tls key file pem format, if not set, the program will automatically generate
  -local-address string
        Local proxy listening address (default "http://0.0.0.0:8080")
  -local-auth string
        Local proxy auth username and password
  -logfile string
        the logging file, using the stdout as default.
  -mode string
        proxy mode(local/proxy/domain/auto) (default "proxy")
  -remote-address string
        Remote proxy listening address (default "https://your.vps:8080")
  -remote-auth string
        Remote proxy auth username and password
  -stat
        display the traffic statistics in a loop timer
  -timeout int
        connect timeout (unit second) (default 30)
```
Of course you can modify and run again

### Everything is ready, return to the home windows and start the service; happy you life;

### If you think this software is good, you can consider asking the author for coffee;

### [paypal.me](https://paypal.me/lixiangyun)