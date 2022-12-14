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
        running mode(local/proxy/domain) (default "proxy")
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
./autoproxy -local-address https://0.0.0.0:8080 -mode local -local-auth user:123456
```

- 本地部署:
```
./autoproxy -local-auth "" -remote-address https://{remote-ip}:8080 -remote-auth user:123456
```

#### 1、本地代理模式；

在本地代理下面，只需要部署一个autoproxy程序，这个程序作为内网主机访问外网的代理服务；配置参考如下：

```
./autoproxy -local-address https://0.0.0.0:8080 -mode local -local-auth user:123456
```

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

#### 2、本地代理+二级代理模式
二级代理就是在本地代理基础之上，将本地代理的部分或者全部流量通过指定二级代理服务进行转发；可用于复杂的网络环境下，部分网站加速；

本地代理配置：

```
./autoproxy -local-auth "" -remote-address https://{remote-ip}:8080 -remote-auth user:123456
```

二级代理部署在VPS侧，需要准备具备一个公网IP的虚拟云主机；

```
./autoproxy -local-address https://0.0.0.0:8080 -mode local -local-auth user:123456
```

#### 3、本地windows UI客户端

本项目提供了小白使用的带UI的客户端，方便使用；在Release 版本下载 `autoproxy_win.zip`
然后解压双击运行即可，本程序是免安装的绿色软件，不会破坏系统；

特性如下：

- 提供基本设置选项
- 转发域名规则
- 远程服务配置
- 最小化和隐藏窗口
- 实时统计控制台
- 本地代理地址和端口设置


![](./docs/main.PNG)

#### 添加二级代理服务

![](./docs/remoteproxy.PNG)

配置完成后，您可以单击“测试”以尝试网络连接性；

#### 编辑域名转发规则

支持几种匹配规则，例如：

- `*.domain.*` : 匹配中间部分域名字段
- `*.domain.com` : 匹配后面域名字段
- `www.domain.*` : 匹配前面域名字段
- `www.domain.com` : 完整匹配域名

![](./docs/domain.PNG)

#### 提供多语言设置

![](./docs/language.PNG)

#### 同步修改本地Internet设置选项

![](./docs/setting.go.PNG)

### 一切就绪，返回主窗口并启动服务； 快乐生活

### [paypal.me](https://paypal.me/lixiangyun)

### Alipay
![](./autoproxy_win/static/sponsor1.jpg)

### Wechat Pay 
![](./autoproxy_win/static/sponsor2.jpg)

### 感谢支持
<a href="https://jb.gg/OpenSource">
<img src="https://github.com/easymesh/autoproxy/blob/master/docs/jetbrains.png" title="Logo" width="100" height="100"/>
</a>