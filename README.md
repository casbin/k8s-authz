# k8s-authz

### 基本概述

​	

### 原理图

![](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/k8s1.png)

### 本插件怎么与Casbin交互

插件引用casbin的库, 在默认端口8888监听 webhook发来的验证请求

接受到后给一个调用casbin的handler处理, 然后返回结果



### K8s怎么与Webhook 通信

- 创建TLS Certificate，即证书
- 编写服务端代码，服务端代码需要使用证书
- 根据证书创建k8s sercret
- 创建k8s Deployment和Service
- 创建k8s WebhookConfiguration，其中需要使用之前创建的证书



