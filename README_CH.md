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





#### 先决条件

确保 Kubernetes 集群版本至少为 v1.16（以便使用 `admissionregistration.k8s.io/v1` API） 或者 v1.9 （以便使用 `admissionregistration.k8s.io/v1beta1` API）

使用本条命令

```
kubectl api-versions | grep admissionregistration.k8s.io
```

结果应该是

```
admissionregistration.k8s.io/v1
admissionregistration.k8s.io/v1beta1
```



#### 测试环境

kubernetes 1.16.7



#### 测试方法

1.配置认证证书

![image-20200913102138806](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200913102138806.png)

2.验证脚本

![image-20200913102324194](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200913102324194.png)

3.部署服务

![image-20200913102412956](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200913102412956.png)

4.替换deployment.yaml文件里的字段,执行完成后可以查看`validatingwebhook-ca-bundle.yaml`文件中的`CA_BUNDLE`占位符的值是否已经被替换掉了。

![image-20200913102502928](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200913102502928.png)

5.在`default`这个 namespace 中添加该标签：

```
$ kubectl label namespace default admission-webhook-example=enabled
namespace "default" labeled
```

![image-20200913102748776](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200913102748776.png)



6.创建一个 deployment 资源来验证下是否有效，代码仓库下有一个`sleep.yaml`的资源清单文件，直接创建即可：

```shell
$ kubectl create -f deployment/sleep.yaml
Error from server (required labels are not set): error when creating "deployment/sleep.yaml": admission webhook "required-labels.qikqiak.com" denied the request: required labels are not set
```

7.正常情况下创建的时候会出现上面的错误信息，然后部署另外一个`sleep-with-labels.yaml`的资源清单：

```shell
kubectl create -f deployment/sleep-with-labels.yaml
deployment.apps "sleep" created
```