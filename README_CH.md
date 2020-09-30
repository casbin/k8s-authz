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



### 代码结构

model里面放的是casbin的model文件
policy里面放的是casbin的policy文件
deployment里面是测试所用到的脚本, yaml文件
casbin_server里面是casbin_server
main.go是启动主函数
webhook.go是主要的权限访问逻辑





### 部署服务

为了部署 webhook server，我们需要在我们的 Kubernetes 集群中创建一个 service 和 deployment  资源对象，部署是非常简单的，只是需要配置下服务的 TLS 配置。我们可以在代码根目录下面的 deployment 文件夹下面查看`deployment.yaml`文件中关于证书的配置声明，会发现从命令行参数中读取的证书和私钥文件是通过一个 secret 对象挂载进来的：

```yaml
args:
	- -tlsCertFile=/etc/webhook/certs/cert.pem
	- -tlsKeyFile=/etc/webhook/certs/key.pem
[...]
	volumeMounts:
	- name: webhook-certs
		mountPath: /etc/webhook/certs
		readOnly: true
volumes:
- name: webhook-certs
  secret:
	secretName: admission-webhook-example-certs
```

在生产环境中，对于 TLS 证书（特别是私钥）的处理是非常重要的，我们可以使用类似于[cert-manager](https://www.qikqiak.com/post/automatic-kubernetes-ingress-https-with-lets-encrypt)之类的工具来自动处理 TLS 证书，或者将私钥密钥存储在Vault中，而不是直接存在 secret 资源对象中。  我们可以使用任何类型的证书，但是需要注意的是我们这里设置的 CA 证书是需要让 apiserver 能够验证的，我们这里可以重用 Istio 项目中的生成的[证书签名请求脚本](https://github.com/istio/istio/blob/release-0.7/install/kubernetes/webhook-create-signed-cert.sh)。通过发送请求到 apiserver，获取认证信息，然后使用获得的结果来创建需要的 secret 对象。

首先，运行[该脚本](https://github.com/cnych/admission-webhook-example/blob/blog/deployment/webhook-create-signed-cert.sh)检查 secret 对象中是否有证书和私钥信息：

```shell
$ ./deployment/webhook-create-signed-cert.sh
creating certs in tmpdir /var/folders/x3/wjy_1z155pdf8jg_jgpmf6kc0000gn/T/tmp.IboFfX97 
Generating RSA private key, 2048 bit long modulus (2 primes)
..................+++++
........+++++
e is 65537 (0x010001)
certificatesigningrequest.certificates.k8s.io/admission-webhook-example-svc.default created
NAME                                    AGE   REQUESTOR          CONDITION
admission-webhook-example-svc.default   1s    kubernetes-admin   Pending
certificatesigningrequest.certificates.k8s.io/admission-webhook-example-svc.default approved
secret/admission-webhook-example-certs created

$ kubectl get secret admission-webhook-example-certs
NAME                              TYPE     DATA   AGE
admission-webhook-example-certs   Opaque   2      28s
```

一旦 secret 对象创建成功，我们就可以直接创建 deployment 和 service 对象。

```shell
$ kubectl create -f deployment/deployment.yaml
deployment.apps "admission-webhook-example-deployment" created

$ kubectl create -f deployment/service.yaml
service "admission-webhook-example-svc" created
```

### 配置 webhook

现在我们的 webhook 服务运行起来了，它可以接收来自 apiserver 的请求。但是我们还需要在 kubernetes 上创建一些配置资源。首先来配置 validating 这个 webhook，查看 [webhook 配置](https://github.com/cnych/admission-webhook-example/blob/blog/deployment/validatingwebhook.yaml)，我们会注意到它里面包含一个`CA_BUNDLE`的占位符：

```yaml
clientConfig:
  service:
	name: admission-webhook-example-svc
	namespace: default
	path: "/validate"
  caBundle: ${CA_BUNDLE}
```

CA 证书应提供给 admission webhook 配置，这样 apiserver 才可以信任 webhook server 提供的  TLS 证书。因为我们上面已经使用 Kubernetes API 签署了证书，所以我们可以使用我们的 kubeconfig 中的 CA  证书来简化操作。代码仓库中也提供了一个小脚本用来替换 CA_BUNDLE 这个占位符，创建 validating webhook  之前运行该命令即可：

```shell
$ cat ./deployment/validatingwebhook.yaml | ./deployment/webhook-patch-ca-bundle.sh > ./deployment/validatingwebhook-ca-bundle.yaml
```

执行完成后可以查看`validatingwebhook-ca-bundle.yaml`文件中的`CA_BUNDLE`占位符的值是否已经被替换掉了。需要注意的是 clientConfig 里面的 path 路径是`/validate`，因为我们代码在是将 validate 和 mutate 集成在一个服务中的。

然后就是需要配置一些 RBAC 规则，我们想在 deployment 或 service 创建时拦截 API 请求，所以`apiGroups`和`apiVersions`对应的值分别为`apps/v1`对应 deployment，`v1`对应 service。对于 RBAC 的配置方法可以查看我们前面的文章：[Kubernetes RBAC 详解](https://www.qikqiak.com/post/use-rbac-in-k8s)

webhook 的最后一部分是配置一个`namespaceSelector`，我们可以为 webhook 工作的命名空间定义一个 selector，这个配置不是必须的，比如我们这里添加了下面的配置：

```yaml
namespaceSelector:
  matchLabels:
	admission-webhook-example: enabled
```

则我们的 webhook 会只适用于设置了`admission-webhook-example=enabled`标签的 namespace， 您可以在Kubernetes参考文档中查看此资源配置的完整布局。

所以，首先需要在`default`这个 namespace 中添加该标签：

```shell
$ kubectl label namespace default admission-webhook-example=enabled
namespace "default" labeled
```

最后，创建这个 validating webhook 配置对象，这会动态地将 webhook 添加到 webhook 链上，所以一旦创建资源，就会拦截请求然后调用我们的 webhook 服务：

```shell
$ kubectl create -f deployment/validatingwebhook-ca-bundle.yaml
validatingwebhookconfiguration.admissionregistration.k8s.io "validation-webhook-example-cfg" created
```

### 测试

现在让我们创建一个 deployment 资源来验证下是否有效，代码仓库下有一个`sleep.yaml`的资源清单文件，直接创建即可：

```shell
$ kubectl create -f deployment/sleep.yaml
Error from server (required labels are not set): error when creating "deployment/sleep.yaml": admission webhook "required-labels.qikqiak.com" denied the request: required labels are not set
```

正常情况下创建的时候会出现上面的错误信息，然后部署另外一个`sleep-with-labels.yaml`的资源清单：

```shell
$ kubectl create -f deployment/sleep-with-labels.yaml
deployment.apps "sleep" created
```

可以看到可以正常部署，先我们将上面的 deployment 删除，然后部署另外一个`sleep-no-validation.yaml`资源清单，该清单中不存在所需的标签，但是配置了`admission-webhook-example.qikqiak.com/validate=false`这样的 annotation，正常也是可以正常创建的：

```shell
$ kubectl delete deployment sleep
$ kubectl create -f deployment/sleep-no-validation.yaml
deployment.apps "sleep" created
```

### 部署 mutating webhook

首先，我们将上面的 validating webhook 删除，防止对 mutating 产生干扰，然后部署新的配置。 mutating webhook 与 validating webhook 配置基本相同，但是 webook server 的路径是`/mutate`，同样的我们也需要先填充上`CA_BUNDLE`这个占位符。

```shell
$ kubectl delete validatingwebhookconfiguration validation-webhook-example-cfg
validatingwebhookconfiguration.admissionregistration.k8s.io "validation-webhook-example-cfg" deleted

$ cat ./deployment/mutatingwebhook.yaml | ./deployment/webhook-patch-ca-bundle.sh > ./deployment/mutatingwebhook-ca-bundle.yaml

$ kubectl create -f deployment/mutatingwebhook-ca-bundle.yaml
mutatingwebhookconfiguration.admissionregistration.k8s.io "mutating-webhook-example-cfg" created
```

现在我们可以再次部署上面的`sleep`应用程序，然后查看是否正确添加 label 标签：

```shell
$ kubectl create -f deployment/sleep.yaml
deployment.apps "sleep" created

$ kubectl get  deploy sleep -o yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    admission-webhook-example.qikqiak.com/status: mutated
    deployment.kubernetes.io/revision: "1"
  creationTimestamp: 2018-09-24T11:35:50Z
  generation: 1
  labels:
    app.kubernetes.io/component: not_available
    app.kubernetes.io/instance: not_available
    app.kubernetes.io/managed-by: not_available
    app.kubernetes.io/name: not_available
    app.kubernetes.io/part-of: not_available
    app.kubernetes.io/version: not_available
...
```

最后，我们重新创建 validating webhook，来一起测试。现在，尝试再次创建 sleep 应用。正常是可以创建成功的，我们可以查看下[admission-controllers 的文档](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#what-are-they)

> 准入控制分两个阶段进行，第一阶段，运行 mutating admission 控制器，第二阶段运行 validating admission 控制器。

所以 mutating webhook 在第一阶段添加上缺失的 labels 标签，然后 validating webhook 在第二阶段就不会拒绝这个 deployment 了，因为标签已经存在了，用`not_available`设置他们的值。

```shell
$ kubectl create -f deployment/validatingwebhook-ca-bundle.yaml
validatingwebhookconfiguration.admissionregistration.k8s.io "validation-webhook-example-cfg" created

$ kubectl create -f deployment/sleep.yaml
deployment.apps "sleep" created
```



参考

https://banzaicloud.com/blog/k8s-admission-webhooks/

https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#what-are-they