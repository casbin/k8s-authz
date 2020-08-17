package main

import (
	"fmt"
	casbin"k8s-authz/casbin_server"
	client"k8s.io/client-go"
)

func main(){

	connectToK8s()
	webhook:=initialize()
	casbin.RunServer()
	listen(webhook)
	for true {

		fmt.Print("hello")
		mes:=getMessageFromWebhook()

		res:=sendToCasbin(mes)
		postToK8s(res)

	}

}
func postToK8s(res string)  {
	
}
func sendToCasbin(mes string)  string{
	return mes+"after process"
}
func  getMessageFromWebhook() string  {
	return "after decode mes"
}

func listen(url string)  {
	//监听k8s的url
}

func initialize() string{
	return "https://k8s_authz.com/authz"
}

func connectToK8s(){

}
