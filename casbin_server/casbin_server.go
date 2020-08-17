package casbin_server

import (
	_"fmt"
	_"github.com/casbin/casbin"
)



func RunServer(){
	listen(6666)
	mes:=getMessage()

	res:=casbin_authz(mes)
	sendMessage(res)

}
func sendMessage(res string)   {
	//发消息回k8s_authz
}
func casbin_authz(mes string) string  {
	//接入Casbin
	//此处使用casbin
	return "after authz"
}
func getMessage() string  {
	return "get the message  from k8s_authz"
}
func listen(port int)  {
	//监听k8s_authz 发来的信息


}
func initialize(){

}

