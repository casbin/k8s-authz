# k8s-authz
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/casbin/k8s-authz/issues)
[![Go Report Card](https://goreportcard.com/badge/github.com/casbin/k8s-authz)](https://goreportcard.com/report/github.com/casbin/k8s-authz)
[![Coverage Status](https://coveralls.io/repos/github/casbin/k8s-authz/badge.svg?branch=master)](https://coveralls.io/github/casbin/k8s-authz?branch=master)
[![Go](https://github.com/casbin/k8s-authz/actions/workflows/ci.yaml/badge.svg)](https://github.com/casbin/k8s-authz/actions/workflows/ci.yaml)
[![Discord](https://img.shields.io/discord/1022748306096537660?logo=discord&label=discord&color=5865F2)](https://discord.gg/S5UjpzGZjN)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

<p align="center">
    <img width="300" height="300" src="k8s-logo.png" alt="K8s-authz" />
</p>
   
K8s-authz is authorization middleware for [Kubernetes](https://github.com/kubernetes/kubernetes), which is based on [Casbin](https://github.com/casbin/casbin). 

## Installation

```
go get github.com/casbin/k8s-authz
```
## Working

This middleware uses K8s validation admission webhook to check the policies defined by casbin, for every request related to the pods. The K8s API server needs to know when to send the incoming request to our admission controller. For this part, we have defined a validation webhook which would proxy the requests for the pods and perform policy verification on it. The user would be allowed to perform the operations on the pods, only if the casbin enforcer authorizes it. The enforcer checks the roles of the user defined in the policies. This middleware would be deployed on the k8s cluster. 

## Requirements
Before proceeding, make sure to have the following-
- Running k8s Cluster
- kubectl
- Openssl

## Configuration and Usage
 
- Generate the certificates and keys for every user by using openssl and running the following script:

  If you are on a Linux system, you can execute shell scripts directly
    ```
    ./gen_cert.sh
    ```
  If you are on a Windows system, executing `./gen_cert.sh` can be problematic, especially if you are using `Git Bash`
  Follow the steps below:
    ```
  # Do not use Git Bash to execute these commands (You can use cmd)
  
    openssl genrsa -out certs/ca.key 2048
    
    openssl req -new -x509 -key certs/ca.key -out certs/ca.crt
    
    openssl genrsa -out certs/casbin-key.pem 2048
    
    openssl req -new -key certs/casbin-key.pem -subj "/CN=casbin.default.svc" -out casbin.csr
    
    openssl x509 -req -in casbin.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/casbin-crt.pem
  
  # You can use Git Bash to execute the following command, or you can use other equivalent methods
    
    export CA_BUNDLE=$(cat certs/ca.crt | base64 | tr -d '\n')
    
    cat manifests/ValidatingWebhookConf.yaml.template | envsubst > manifests/ValidatingWebhookConf.yaml
    ```

- For a production server, we need to create a k8s `secret` to place the certificates for security purposes.
    ```
    kubectl create secret generic authz -n default \
      --from-file=key.pem=certs/casbin-key.pem \
      --from-file=cert.pem=certs/casbin-crt.pem
    ```
- Once, this part is done we need to change the directory of the certs in [main.go](https://github.com/ashish493/k8s-authz/blob/3560551427c0431a9d4594ad1206f084ede37c49/main.go#L26) and then in [manifests](https://github.com/ashish493/k8s-authz/blob/3560551427c0431a9d4594ad1206f084ede37c49/manifests/deployment.yaml#L22) with that of the `secret`.

- Build the docker image from the [Dockerfile](https://github.com/casbin/k8s-authz/blob/master/Dockerfile) manually by running the following command and then change the build version here and at the deployment [file](https://github.com/casbin/k8s-authz/blob/718f58c46e3dbf79063b5b1c18348c2fee5de9e9/manifests/deployment.yaml#L18), as per the builds.
    ```
    docker build -t casbin/k8s_authz:latest .
    ```
  
- Define the casbin policies in the [model.conf](https://github.com/casbin/k8s-authz/blob/master/config/model.conf) and [policy.csv](https://github.com/casbin/k8s-authz/blob/master/config/policy.csv). You can refer the [docs](https://casbin.org/docs/how-it-works) to get to know more about the working of these policies.

- Before deploying, you can change the ports in [main.go](https://github.com/casbin/k8s-authz/blob/master/main.go) and also in the validation webhook configuration [file](https://github.com/casbin/k8s-authz/blob/master/manifests/deployment.yaml) depending on your usage.

- Deploy the validation controller and the webhook on k8s cluster by running:
    ```
    kubectl apply -f manifests/deployment.yaml
  
    # Wait for Deployment Ready
  
    kubectl apply -f manifests/ValidatingWebhookConf.yaml
  ```

Now the server should be running and ready to validate the requests for the operations on the pods. 

## Documentation

You can check the official [docs](https://casbin.org/docs/k8s) for more detailed explanation.

## Community

In case of any query, you can ask on our [Discord](https://discord.gg/S5UjpzGZjN).

