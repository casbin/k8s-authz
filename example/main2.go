package main

import (
	"fmt"
	"github.com/casbin/casbin"
	_"k8s-authz/model"

)

func main()  {
	e := casbin.NewEnforcer("k8s-authz/model/auth_model.conf", "policy/policy.csv")
	sub := "alice" // the user that wants to access a resource.
	obj := "data1" // the resource that is going to be accessed.
	act := "read" // the operation that the user performs on the resource.
	if res := e.Enforce(sub, obj, act); res {
		// permit alice to read data1
		fmt.Println(res)
	} else {
		fmt.Println("error")
		// deny the request, show an error
	}
}
