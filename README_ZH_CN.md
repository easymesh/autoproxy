# Autoproxy

[English](./README.md)
[中文](./README_ZH_CN.md) 

该项目提供浏览器http proxy代理服务，支持https、http协议代理；可以部署多级代理；支持TLS协议加密；主要使用与内网主机访问外部网站工具；

## 特性如下
- 支持多种转发方式
- 支持统计控制台
- 支持二次转发
- 支持TLS传输加密
- 支持动态路由
- 支持自定义转发域名设置
- 支持多种平台

## 部署
### 准备提前
- 准备具有公共IP的虚拟云主机
- 准备指定外部TCP协议访问端口，例如：8080

### 使用方式：
下载解压相应平台的软件包；其中包括三个文件；一个autoproxy可执行程序;
程序命令参数为：

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
        running mode(local/proxy/domain/auto) (default "proxy")
  -remote-address string
        Remote proxy listening address (default "https://you.domain.com:8080")
  -remote-auth string
        Remote proxy auth username and password (default "user:passwd")
  -timeout int
        connect timeout (unit second) (default 30)
```

后台启动参考：
- 服务部署:
```
./autoproxy -local-address https://0.0.0.0:8080 -mode local -local-auth user:password
```

- 本地部署:
```
./autoproxy -local-auth "" -remote-address https://{remote-ip}:8080 -remote-auth user:password
```

#### 1、本地代理模式；

在本地代理下面，只需要部署一个autoproxy程序，这个程序作为内网主机访问外网的代理服务；配置参考如下：

```
./autoproxy -local-address https://0.0.0.0:8080 -mode local -local-auth user:password
```

然后设置浏览器或者环境变量；
```
export http_proxy="http://用户名:密码@一级代理IP:端口"
export https_proxy="http://用户名:密码@一级代理IP:端口"
```

例如：
```
export http_proxy="http://user:password1@192.168.3.1:8080"
export https_proxy="http://user:password1@192.168.3.1:8080"
```

#### 2、本地代理+二级代理模式
二级代理就是在本地代理基础之上，将本地代理的部分或者全部流量通过指定二级代理服务进行转发；可用于复杂的网络环境下，部分网站加速；

本地代理配置：

```
./autoproxy -local-auth "" -remote-address https://{remote-ip}:8080 -remote-auth user:password
```

二级代理部署在VPS侧，需要准备具备一个公网IP的虚拟云主机；

```
./autoproxy -local-address https://0.0.0.0:8080 -mode local -local-auth user:password
```

### 一切就绪，返回主窗口并启动服务； 快乐生活

### [paypal.me](https://paypal.me/lixiangyun)