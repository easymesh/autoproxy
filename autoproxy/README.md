# 【Auto Proxy】自研Go语言自动代理工具软件

该项目提供浏览器http proxy代理服务，支持https、http协议代理；可以部署多级代理；支持TLS协议加密；主要使用与内网主机访问外部网站工具；

当前版本特性如下：

1、支持TLS传输加密；
2、支持动态路由；
3、支持多账户认证；
4、支持日志异步记录；
5、支持多系统平台;

### 官方发布：

[https://www.easymesh.info/](https://www.easymesh.info/)

### 版本下载：

[https://github.com/easymesh/autoproxy/releases](https://github.com/easymesh/autoproxy/releases)

### 使用方式：

1、下载解压相应平台的软件包；其中包括三个文件；一个autoproxy可执行程序，以及两个参考配置文件；server.yaml 和 client.yaml;

程序命令参数为：

```
Usage of autoproxy.exe:
  -config string
        configure file (default "config.yaml")
  -debug
        enable debug
  -help
        usage help
```

`autoproxy -config server.yaml 或者 autoproxy -config client.yaml`
2、一级代理参考；

在一级代理下面，只需要部署一个autoproxy程序，这个程序作为内网主机访问外网的代理服务；配置参考如下：

```
log:
  path: ./
  filesize: 10485760
  filenumber: 60
local:
  listen: 0.0.0.0:8080
  timeout: 30
  auth:
    - username: user1
      password: uS31k5KLh3NyfvHtFk
    - username: user2
      password: c2O9XJGG0bsJLpt6tr
  mode: local
```

- log 表示日志记录的目录、单个文件大小、以及文件数量上限；主要是用于审计和问题定位；
- local 表示程序提供的服务配置，包括监听地址和端口，链路超时时间以及认证方式；如果没有配置认证；则不会进行认证；
然后设置浏览器或者环境变量；

```
export http_proxy="http://用户名:密码@一级代理IP:端口"
export https_proxy="http://用户名:密码@一级代理IP:端口"
```

例如：
```
export http_proxy="http://user1:password1@192.168.3.1:8080"
export https_proxy="http://user1:password1@192.168.3.1:8080"
```

3、二级代理参考；

二级代理就是在一级代理基础之上，将一级代理的部分或者全部流量通过指定二级代理服务进行转发；可用于复杂的网络环境下，部分网站加速；

一级代理配置；

参考压缩包的client.yaml配置文件，默认只需要修改指定二级代理IP地址就可以使用了；

```
log:
  path: ./
  filesize: 10485760
  filenumber: 60
local:
  listen: 0.0.0.0:8080
  timeout: 5
  mode: auto
remote:
  - address: {二级代理IP}:8080
    timeout: 30
    auth:
      username: user1
      password: uS31k5KLh3NyfvHtFk
    tls:
      enable: true
```
- local: 其中 mode 有三个选项，分别是：local、auto、proxy ，其中local 表示所有流量通过本地路由处理，不会经过二级代理；auto 表示根据IP可达性，比如有些本地路由访问不了或者链路超时，则会使用二级代理进行转发，proxy 表示所有流量全部经过二级代理；
- remote: 需要访问一个或者多个二级代理的地址，超时时间，认证信息；是否进行TLS加密；如果配置多个二级地址，那么会逐个进行链接尝试；

二级代理配置：
```
log:
  path: ./
  filesize: 10485760
  filenumber: 60
local:
  listen: 0.0.0.0:8080
  timeout: 30
  auth:
    - username: user1
      password: uS31k5KLh3NyfvHtFk
    - username: user2
      password: c2O9XJGG0bsJLpt6tr
  mode: local
  tls:
    enable: true
```
改配置表示二级代理服务端口、认证信息，是否进行TLS加密；如果未配置TLS加密传输，那么一级代理的remote的TLS配置也需要去掉；否则就会链接失败；

 

声明：该工具作为免费软件授权使用，软件著作权归作者所有，使用和传播必须符合国内法律法规，如果违反任何法律法规与本人无关；本人对于任何原因在使用本软件对用户自己或者他人造成的任何形式的损失和伤害不承担任何责任；
