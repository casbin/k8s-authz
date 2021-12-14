FROM golang:1.15-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN mkdir ~/k8s-authz
WORKDIR ~/k8s-authz  

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o ./authz

EXPOSE 443
CMD ["./authz"]