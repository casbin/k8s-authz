# 基于 Kubernetes 构建云原生分布式访问控制应用

# 中期报告

## 项目信息

- 项目名称：基于 Kubernetes 构建云原生分布式访问控制应用

- 方案描述： 写一个插件用来作为与K8s和Casbin通信的中间件. 利用Kubernetes提供的client-go 和kubernetes集群的master节点进行通信, 我们再启动另外一个进程用来跑Casbin,其中与通信的协议需要自己定. 与kubernetes通信比较复杂 , 我们首先得获取到master节点上的config文件, 这部分代码已经有了. 然后就是

  k8s需要开启Webhook功能

  其中开启webhook的代码

  ```
  --admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota
  ```

  [开启webhook的参数ValidatingAdmissionWebhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/)

  > ## 一、Webhook概念
  
  >Webhook是一个API概念，
  >
  >**webhook机制**:项目A提供一个webhook url,每次项目B创建新数据时,便会向项目A的hook地址进行请求,项目A收到项目B的请求,然后对数据进行处理
  >
  >准确的说webhook是一种web回调或者http的push API，是向APP或者其他应用提供实时信息的一种方式。Webhook在数据产生时立即发送数据，也就是你能实时收到数据。这一种不同于典型的API，需要用了实时性需要足够快的轮询。这无论是对生产还是对消费者都是高效的，唯一的缺点是初始建立困难。Webhook有时也被称为反向API，因为他提供了API规则，你需要设计要使用的API。Webhook将向你的应用发起http请求，典型的是post请求，应用程序由请求驱动。
  >
  >## 二、使用webhook
>
  >消费一个webhook是为webhook准备一个URL，用于webhook发送请求。这些通常由后台页面和或者API完成。这就意味你的应用要设置一个通过公网可以访问的URL。
>
  >多数webhook以两种数据格式发布数据：JSON或者XML，这需要解释。另一种数据格式是application/x-www-form-urlencoded or multipart/form-data。这两种方式都很容易解析，并且多数的Web应用架构都可以做这部分工作。
  >
  >**项目A需要实时获取到项目B的最新数据：**
  >
  >传统做法:项目A需要不停轮询去拉取项目B的最新数据
  >
  >![img](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/1353331-20200617135512464-1320756712.png)

  在K8s开启webhook后 我们就可以通过这个和APIserver进行通信了

  原理也就是 当Casbin接受到一个访问请求时, 会进行权限验证, 把处理后的结果 , 整理成json格式, 然后向 APIserve的url 发送一个POST请求, POST请求里会权限验证的最后结果,通知给APIserver, 然后APIserver就以这个数据来

  



​		然后就是Casbin进程跑在哪里, 我们把Casbin进程当成黑盒用, 目前计划的是把Casbin打包成镜像 让这个docker镜像跑在集群里面

这样有两个好处

1. 可以通过改镜像标签来动态更新Casbin的版本

2. 可以让Casbin跑在集群里 , 不用再额外管理Casbin 进程, 比如可以让Casbin 实例数为3个, 这样可以同时有3个Casbin可以处理访问控制请求, 也可以随时进行扩容

   其中大概是这样子的

   ```
   containers:
           - name: Casbin
             image: Casbin/Casbin:latest
             args:
               - "run"
               - "--server"
               - "--tls-cert-file=/certs/tls.crt"
               - "--tls-private-key-file=/certs/tls.key"
               - "--addr=0.0.0.0:443"
               - "--addr=http://127.0.0.1:8181"
               - "--log-format=json-pretty"
               - "--set=decision_logs.console=true"
             volumeMounts:
               - readOnly: true
                 mountPath: /certs
                 name: opa-server
             readinessProbe:
               httpGet:
                 path: /health?plugins&bundle
                 scheme: HTTPS
                 port: 443
               initialDelaySeconds: 3
               periodSeconds: 5
             livenessProbe:
               httpGet:
                 path: /health
                 scheme: HTTPS
                 port: 443
               initialDelaySeconds: 3
               periodSeconds: 5
           - name: kube-mgmt
             image: casbin/kube-mgmt:0.11
             args:
               - "--replicate-cluster=v1/namespaces"
               - "--replicate=extensions/v1beta1/ingresses"
         volumes:
           - name: casbin
             secret:
               secretName: casbin-server
   ```

   

- 时间规划：

  # Week 1 (June 27 - July 5)

  ## Weekly Summary

  1. Studied the fundamental knowledge about Golang(Actually I am moving  from C++ to Golang, so I am not very familiar with everything in  Golang).
  2. Still setting up the environment (including the k8s and Ubuntu for service test).
  3. Conducted further investigation on other authorization library.
  4. Read the paper 

  ## Before Coding

  Before get started, I think we should have  further knowledge  of   distributed system, so I go to learn the MIT-6.824 . Only after I have a solid foundation can I go further.

  # Week 2 (June 6 - July 12)

  ## Weekly Summary

  1. Studied the fundamental knowledge about Kubenetes, including Pods, Service ,Kubenetes APIServer ,Kubenetes Object , Node, ReplicaSet
  2. build up two centos 7.8 server

  ### k8s 环境配置

  1.改用户名
  检查 centos / hostname

  ```
  # 在 master 节点和 worker 节点都要执行
  cat /etc/redhat-release
  
  # 此处 hostname 的输出将会是该机器在 Kubernetes 集群中的节点名字
  # 不能使用 localhost 作为节点的名字
  hostname
  
  修改 hostname
  
  如果您需要修改 hostname，可执行如下指令：
  
  # 修改 hostname
  hostnamectl set-hostname your-new-host-name
  # 查看修改结果
  hostnamectl status
  # 设置 hostname 解析
  echo "127.0.0.1   $(hostname)" >> /etc/hosts
  # 请使用 lscpu 命令，核对 CPU 信息
  # Architecture: x86_64    本安装文档不支持 arm 架构
  # CPU(s):       2         CPU 内核数量不能低于 2
  lscpu
  ```

  2.配置网络

  ```
  1.首先dhclient 通一下网络
  2.vim /etc/sysconfig/network-scripts-ifcfg-ens33
  3.systemctl restart network.service
  
  宿主机ip
  10.2.32.100
  掩码
  255.255.248.0
  网关
  10.2.32.1
  
  master结点
  ip
  10.2.34.9
  
  worker节点
  ip
  10.2.34.73
  ```

  3.安装docker    docker nfs-utils kubectl / kubeadm / kubelet

  ```
  # 在 master 节点和 worker 节点都要执行
  # 最后一个参数 1.18.5 用于指定 kubenetes 版本，支持所有 1.18.x 版本的安装
  # 腾讯云 docker hub 镜像
  # export REGISTRY_MIRROR="https://mirror.ccs.tencentyun.com"
  # DaoCloud 镜像
  # export REGISTRY_MIRROR="http://f1361db2.m.daocloud.io"
  # 华为云镜像
  # export REGISTRY_MIRROR="https://05f073ad3c0010ea0f4bc00b7105ec20.mirror.swr.myhuaweicloud.com"
  # 阿里云 docker hub 镜像
  export REGISTRY_MIRROR=https://registry.cn-hangzhou.aliyuncs.com
  curl -sSL https://kuboard.cn/install-script/v1.18.x/install_kubelet.sh | sh -s 1.18.5
  
  ```

  4.初始化master

  ```
  初始化 master 节点
  
  关于初始化时用到的环境变量
  
      APISERVER_NAME 不能是 master 的 hostname
      APISERVER_NAME 必须全为小写字母、数字、小数点，不能包含减号
      POD_SUBNET 所使用的网段不能与 master节点/worker节点 所在的网段重叠。该字段的取值为一个 CIDR 值，如果您对 CIDR 这个概念还不熟悉，请仍然执行 export POD_SUBNET=10.100.0.1/16 命令，不做修改
  
  
  
  请将脚本最后的 1.18.5 替换成您需要的版本号， 脚本中间的 v1.18.x 不要替换
  ```

  ```
  # 只在 master 节点执行
  # 替换 x.x.x.x 为 master 节点实际 IP（请使用内网 IP）
  # export 命令只在当前 shell 会话中有效，开启新的 shell 窗口后，如果要继续安装过程，请重新执行此处的 export 命令
  export MASTER_IP=x.x.x.x
  # 替换 apiserver.demo 为 您想要的 dnsName
  export APISERVER_NAME=apiserver.demo
  # Kubernetes 容器组所在的网段，该网段安装完成后，由 kubernetes 创建，事先并不存在于您的物理网络中
  export POD_SUBNET=10.100.0.1/16
  echo "${MASTER_IP}    ${APISERVER_NAME}" >> /etc/hosts
  curl -sSL https://kuboard.cn/install-script/v1.18.x/init_master.sh | sh -s 1.18.5
  
   
  
  ```

  5.初始化worker

  6.将worker join

  # Week 3 (July 13- July 12)

  ## Weekly Summary

  1. Studied the fundamental knowledge about kubenetes
  2. build up a kubenetes cluster with single master node and worker node
  3. deploy nginx service on the cluster and aPublish, scale, expand, roll update applications.
  4. load balancing achieve by kubenetes services

  # Week 4 (July 20 - July 27)

  ## Weekly Summary

  ### 7/21

  1.k8s的网络模型

  2.跑通了github.com/client-go/ 的example 利用client-go 与apiserver进行通信

  

  ### 7/22

  1.开始设计以Role为核心的控制访问

  2.学习rego query langanuage

  

  ### 7/23

  1.学习Kubernetes 的Kubernetes Admission Controllers的两个中的Mutating Webhook

  2.
  要想看懂OPA 首先得把OPA

  玩熟练了

  相关的文档读了

  才看得到他代码的思路

  现在OPA好像是 做了个gatekeeper

  然后plain opa就是像你说的黑盒,

  gatekeeper和我们要做的东西重叠度很高

  然后gatekeeper的原理大概就是 建一个mutating webhook对象 通过这个来拦截 直接kubectl直接到apiserver的 请求, 并给予修改, 来决定是否改成deny

  3.使用patch方式来修改






​	2020/8/15-2020/8/18

​		将集群配到云服器上

​	2020/8/18-2020/8/20

​		把go-client的代码再好好读读

​	2020/8/20- 2020/8/30

​		研究研究OPA的代码实现

​	

## 项目进度

- 已完成工作：
  
- 主要完成的工作是前期调研, 架构设计和demo尝试
  
  现在已经可以和k8s APIserver通信了, 我们可以通过client-go提供的API来查询 集群的节点数量 , 动态创建或者
  
  部署Casbin到集群里也应该 没有问题
  
  下一步主要是 定一下通信协议 还有通信过程的代码 这里涉及到大量端口转发 端口的过程.......
  

### 前期调研

由于项目是从0到1的过程,所以在以下几个方案中

1. 研究怎么与K8s通信, 实现一个client,将casbin的功能在这个client中实现.
2. 写一个中间件,将casbin,中间件,kubernetes集群分离开, 实现了低耦合.也就是我们现在采用的方法



方案1的架构图我也设计了

![image-20200815154727987](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200815154727987.png)

选择方案2的原因

1.是可以解耦合, casbin的版本更新可以通过镜像标签的更新来完成,一分钟之内就可以做到, 方便系统的升级

2.选择方案1的话, casbin的功能更新要在这个client里面维护, 后期维护工作量大



	###  架构设计



​		在介绍架构设计之前,首先介绍一下架构里主要的一些组件

1.首先K8s里面Pod概念, Pod是K8s里面最基本的单位,Pod里面跑了镜像的一个实例,比如casbin

​	pod有不同的组件 有些还带了数据卷的,像这样,紫色那个就是数据卷,我们的casbin跑在绿色的正方形里

![image-20200815153038726](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200815153038726.png)

![image-20200815150912649](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200815150912649.png)



2.结点

结点里可以同时跑多个Pod, 我们可以在每一个结点上运行一个casbin应用, 这样本地验证会快一些,也可以集中式的让一个node承担验证的工作

![image-20200815151059626](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200815151059626.png)



3.部署的节点应该是这样

![image-20200815151434080](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200815151434080.png)

黄色的serviceA里面跑的是casbin

绿色的serviceB里面跑的我们的应用,

casbin通过A这个Deployment.yaml文件来声明

应用通过B这个Deployment.yaml来声明



4,最终设计的架构图

![k8s-authz架构图](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/k8s-authz%E6%9E%B6%E6%9E%84%E5%9B%BE.png)

### demo尝试

部署nginx到k8s集群上, 并且用sevice实现负载均衡

1.两个节点,一个master,一个worker的kubernetes集群



![](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/FD6636CD938959D68CAEA132F8BD6EB4.png)



![image-20200815154928024](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200815154928024.png)

![QQ截图20200815155916](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/QQ%E6%88%AA%E5%9B%BE20200815155916.png)2.yaml文件放置demo里了





- ![nginx-demo](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/nginx-demo.png)
- ![nginx-service.demo](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/nginx-service.demo.png)
- ![kubeadm-config](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/kubeadm-config.png)







- 遇到的问题及解决方案：
  
  遇到的问题还是许多环境配置的问题 K8s真的太复杂了...
  
  主要还是有一个体会, 作为学生, 一定要基础扎实才能啃这些, 涉及到的知识是学校教不了的,  但是没有学校教的知识也很难走的远, 有种互补的感觉
  
  
  
- 后续工作安排：
  1.后续会重新部署集群到云服务器上 , 原来配的虚拟机环境 在家里, 现在在学校已经用不了了
  
  2.会定义一种简明的通信API,  完成插件的后续编码工作

