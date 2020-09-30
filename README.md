# k8s-authz

### Basic overview

​	

### Schematic diagram

![](https://raw.githubusercontent.com/yahoo17/MarkdownPictureRepository/master/img/k8s1.png)

### How does this plug-in interact with casbin

The plug-in refers to the library of casbin and listens for the verification request from webhook on the default port 8888
After receiving, it will be processed by a handler calling caspin, and the result will be returned



### How does k8s communicate with webhook

-Create a TLS certificate, that is, a certificate
-Write the server-side code, the server-side code needs to use the certificate
-Create k8s sercret based on certificate
-Create k8s deployment and service
-Create k8s webhookconfiguration, where you need to use the certificate you created earlier





#### Precondition

Ensure that the kubernetes cluster version is at least v1.16（In order to use `admissionregistration.k8s.io/v1` API） or v1.9 （In order to use `admissionregistration.k8s.io/v1beta1` API）

Use this command to check

```
kubectl api-versions | grep admissionregistration.k8s.io
```

The result should be

```
admissionregistration.k8s.io/v1
admissionregistration.k8s.io/v1beta1
```



#### Testing environment

kubernetes 1.16.7