# 基于 Kubernetes 构建云原生分布式访问控制应用

# 中期报告

## 项目信息

- 项目名称：基于 Kubernetes 构建云原生分布式访问控制应用

- 方案描述： 利用Kubernetes提供的client-go将Casbin 进行插件化改造，对Kubernetes的请求（比如pod资源创建等）进行权限约束。

  ####  架构设计

  1.结点里可以同时跑多个Pod， 我们可以在每一个结点上运行一个casbin应用， 这样本地验证会快一些，也可以集中式的让一个node承担验证的工作

  ![image-20200815151059626](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200815151059626.png)

  

  2.部署的节点应该是这样(这种就是集中式的部署)

  ![image-20200815151434080](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200815151434080.png)

  黄色的serviceA里面跑的是casbin

  绿色的serviceB里面跑的我们的应用，

  casbin通过A这个Deployment.yaml文件来声明

  应用通过B这个Deployment.yaml来声明




3.最终设计的架构图

![k8s1](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/k8s1.png)

### 核心用例

​	名称: 对应用Pod创建的管理

​	例子: 比如用户A要申请创建一个nginx的应用



1.首先APIsever 会通知kubernetes插件(也就是我们的项目K8s-authz) , APIserver通过一个url来访问插件, 比如说图中的https://authz.k8s-authz.com/authorize 

2.然后插件会向casbin server发送请求, casbin来判断用户A创建nginx的应用是否合法,假设是不合法的

3.然后casbin server返回结果给k8s-authz, k8s-authz再把结果返回给APIserver

4.APIserver接受到消息, 查看后是不合法的, 于是拒绝请求







### 时间规划：

第1周（6月27日至7月5日）
1.研究了Golang的基本知识（实际上我正在从C++到Golang，所以我对Golang的一切都不太熟悉）。
2.配置环境（包括用于服务测试的k8s和Centos）。
3.对其他User authentication库进行了进一步调查。
4.读Casbin的有关论文， 基于元模型的策略等
5.编码前，我想我们应该对分布式系统有更深入的了解，所以我去学习MIT-6.824。只有我有了坚实的基础，我才能走得更远。
第2周（6月6日至7月12日）
1.学习了Kubenetes的基本知识，包括Pods、服务、Kubenetes APIServer、Kubenetes对象、节点、复制集
2.构建两个centos 7.8服务器
第3周（7月13日-7月22日）
1.学习kubernetes的基本知识
2.建立一个只有主节点和工作节点的kubenetes集群
3.在集群上部署nginx服务，发布、扩展、滚动更新应用。
4.kubenetes services实现的负载平衡

 第4周(7月22日-7月29日)

 1.k8s的网络模型

 2.跑通了github.com/client-go/ 的example 利用client-go 与apiserver进行通信

 第5周(7月30日-8月5日)

 1.开始设计以Role为核心的控制访问

 2.学习rego query langanuage

 第6周(8月5日-8月12日)

 1.学习Kubernetes 的Kubernetes Admission Controllers的两个中的Mutating Webhook

 2.确定架构设计大概就是 建一个mutating webhook对象 通过这个来拦截 直接kubectl直接到apiserver的 请求， 并给予修改， 来决定是否改成deny

 3.使用patch方式来修改

 第7周(8月12日-8月19日)

 1.将集群配到云服器上

 2.完成控制Pod生成的销毁的功能

第8周(8月19日-8月26日)

1.再深入理解go-client的代码

2.研究OPA的代码实现， 并调试Webhook

第9周(8月19日-8月26日)

完成通过APIserver进行资源控制的完整功能， 也就是go-client端

第10周(8月27日-9月3日)

将Casbin打包成镜像， 固定好端口号， 编写与集群里casbin通信的demo以及协议

第11周(9月3日-9月10日)

完成对casbin端的编码，测试，部署

第12周(9月11日-9月18日)

完成插件的代码， 并对client-go和casbin的连接进行调试 ，验证最终功能

第13周(9月19日-9月26日)

完成k8s-authz的上线测试， 并根据反馈进行一些修改

第14周(9月26日-9月30日)

完成K8s-authz项目，准备终期考核

​	

## 项目进度

- 已完成工作：
  
- 主要完成的工作是前期调研， 架构设计和demo尝试
  
  
  

### 前期调研

由于项目是从0到1的过程，所以在以下几个方案中

1. 研究怎么与K8s通信， 实现一个client，将casbin的功能在这个client中实现.
2. 写一个中间件，将casbin，中间件，kubernetes集群分离开， 实现了低耦合.也就是我们现在采用的方法



方案1的架构图我也设计了

![image-20200815154727987](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200815154727987.png)

选择方案2的原因

1.是可以解耦合， casbin的版本更新可以通过镜像标签的更新来完成， 方便系统的升级

2.选择方案1的话， casbin的功能更新要在这个client里面维护， 不方便维护

![k8s1](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/k8s1.png)





### demo尝试

现在已经可以和k8s APIserver通信了， 我们可以通过client-go提供的API来查询 集群的节点数量 ， 动态创建或者

部署Casbin到集群里也应该 没有问题



我们用nginx来模拟我们部署的应用， 通过对nginx节点的创建删除来模拟经过casbin验证后对Pod资源的管理



1.两个节点，一个master，一个worker的kubernetes集群

![](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/FD6636CD938959D68CAEA132F8BD6EB4.png)



![image-20200815154928024](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/image-20200815154928024.png)

![QQ截图20200815155916](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/QQ截图20200815155916.png)

- 2.部署nginx的 yaml截图放置demo里

  

  

  - ![nginx-demo](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/nginx-demo.png)
  - ![nginx-service.demo](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/nginx-service.demo.png)
  - ![kubeadm-config](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/kubeadm-config.png)

  

  

  

  

  

- 遇到的问题及解决方案：

  遇到的问题还是许多环境配置的问题 K8s的实现非常复杂

  主要还是有一个体会， 作为学生， 一定要基础扎实才能啃这些， 涉及到的知识是学校教不了的，  但是没有学校教的知识也很难走的远， 有种互补的感觉

  

- 后续工作安排：
  1.后续会重新部署集群到云服务器上 ， 原来配的虚拟机环境 在家里， 现在在学校已经用不了了

  2.会定义一种简明的通信API，  完成插件的后续编码工作

  3.编写通信过程的代码 把client-go和casbin都接入代码，后续会重新部署集群到云服务器上， 并进行线上的测试