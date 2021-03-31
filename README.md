# Autoproxy

[English](./README.md)
[中文](./README_ZH_CN.md) 

The project provides browser http proxy proxy service, supports https and http protocol proxy; can deploy multi-level proxy; supports TLS protocol encryption; mainly uses tools for accessing external websites with intranet hosts;

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
- Run `tar -zxf autoproxy_linux_amd64.tar.gz` Unzip the compressed package.
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
  -debug
        enable debug
  -domain string
        match domain list file(domain mode requires) (default "domain.json")
  -help
        usage help
  -local-address string
        Local proxy listening address (default "http://0.0.0.0:8080")
  -local-auth string
        Local proxy auth username and password (default "user:passwd")
  -logdir string
        runlog path (default "./")
  -mode string
        running mode(local/proxy/domain) (default "proxy")
  -remote-address string
        Remote proxy listening address (default "https://you.domain.com:8080")
  -remote-auth string
        Remote proxy auth username and password (default "user:passwd")
  -timeout int
        connect timeout (unit second) (default 30)
```
Of course you can modify and run again

### Download and run the windows client
- Choose the latest windows desktop client， such as. `autoproxy_desktop.zip`
- Unzip and double-click to run
- Start successfully, you need to manually add remote proxy service.

#### Home windows
- Provide basic setting options
- Forwarding domain name rules
- Remote service configuration
- Minimize and hide the window
- Real-time statistics console
- Local proxy address and port settings

![](./docs/main.PNG)

#### Add remote service

![](./docs/remoteproxy.PNG)

After the configuration is complete, you can click Test to try to test the connectivity

#### Edit forward domain

Support fuzzy matching rules, For example:

- `*.domain.*` : Middle field matching
- `*.domain.com` : Suffix match
- `www.domain.*` : Prefix match
- `www.domain.com` : Exact match

![](./docs/domain.PNG)


#### Multi-language support

![](./docs/language.PNG)

#### Synchronously modify local Internet setting options

![](./docs/setting.go.PNG)

### Everything is ready, return to the home windows and start the service; happy you life;

### If you think this software is good, you can consider asking the author for coffee;

### [paypal.me](https://paypal.me/lixiangyun)

### Alipay
![](./autoproxy_win/static/sponsor1.jpg)

### Wechat Pay 
![](./autoproxy_win/static/sponsor2.jpg)

### Thanks Support
<a href="https://jb.gg/OpenSource">
<img src="https://github.com/easymesh/autoproxy/blob/master/docs/jetbrains.png" title="Logo" width="100" height="100"/>
</a>