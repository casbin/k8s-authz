package connect

import (
	"log"
	"encoding/json"
	"fmt"
	_"encoding/json"
)

type k8sRequest struct{
	Apiserver string `json:"apisever"`
	Name string	`json:"name"`
	Kind string	`json:"kind"`
	Metadata string `json:"metadata"`
	Number int `json:"number,omitempty"`
	//只有大写的才会被编入
	//omitempty是  omit是忽略的意思
	// empty是空 所以是为0 时忽略掉
}
func Encode(src string){

	
}
func Decode() {
	a:=k8sRequest{
		Apiserver:"v1",
		Name:"webhook", 
		Kind:"mutating webhook",
		Metadata :"coke",
		Number: 0,
	}
	rs,err:= json.Marshal(a)
	if err!=nil{
		log.Fatalln(err)
	}
	fmt.Println("hello json")
	fmt.Println(rs)
	fmt.Println(string(rs))
	
}
