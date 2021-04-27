#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

openssl genrsa -out certs/ca.key 2048

openssl req -new -x509 -key certs/ca.key -out certs/ca.crt 

openssl genrsa -out certs/casbin-key.pem 2048

openssl req -new -key certs/casbin-key.pem -subj "/CN=casbin.default.svc" -out casbin.csr 

openssl x509 -req -in casbin.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/casbin-crt.pem

export CA_BUNDLE=$(cat certs/ca.crt | base64 | tr -d '\n')
cat webhook.yaml | envsubst > _webhook.yaml 