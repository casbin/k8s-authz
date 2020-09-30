package casbin_server
import (
	"fmt"
	"github.com/casbin/casbin"
	"log"
	"net/http"
	"strconv"

)
//The global variables created in this file are referenced by webhook
var AuthEnforcer = casbin.NewEnforcer("../model/auth_model.conf", "policy/auth_policy.csv")
var AclEnforcer = casbin.NewEnforcer("../model/acl_model.conf","policy/acl_policy.csv")
var RestfulEnforcer = casbin.NewEnforcer("../model/restful_model.conf","policy/restful_policy.csv")

//You can also run caspin alone as a permission authentication server
func runCasbin_server()  {
	mux:=http.NewServeMux()
	mux.HandleFunc("/restful",restful)
	mux.HandleFunc("/auth",auth)
	mux.HandleFunc("/acl",acl)
	log.Fatal(http.ListenAndServe("localhost:9000",mux))

}


func restful(w http.ResponseWriter, r * http.Request){
	if err:= r.ParseForm();err!=nil{
		log.Print(err)
	}
	sub:=r.Form["sub"]
	obj:=r.Form["obj"]
	act:=r.Form["act"]
	if res:=RestfulEnforcer.Enforce(sub[0],obj[0],act[0]);res{
		fmt.Fprintf(w,sub[0]+" "+obj[0]+" "+act[0]+" "+strconv.FormatBool(res))
	}else {
		fmt.Fprintf(w,"验证失败")
	}



}
func auth (w http.ResponseWriter,r * http.Request){
	if err:=r.ParseForm();err!=nil{
		log.Print(err)
	}
	sub := r.Form["sub"]
	obj:=r.Form["obj"]
	act:=r.Form["act"]
	if res := AuthEnforcer.Enforce(sub[0], obj[0], act[0]); res {
		// permit alice to read data1
		fmt.Fprintf(w,sub[0]+" "+obj[0]+" "+act[0]+" "+strconv.FormatBool(res))
	} else {
		fmt.Fprintf(w,		strconv.FormatBool(res))

		// deny the request, show an error
	}

}
func acl(w http.ResponseWriter , r * http.Request)  {
	if err:=r.ParseForm();err!=nil{
		log.Print(err)
	}
	sub:=r.Form["sub"]
	obj:=r.Form["obj"]
	act:=r.Form["act"]
	if res:=AclEnforcer.Enforce(sub[0],obj[0],act[0]);res{
		fmt.Fprintf(w,sub[0]+" "+obj[0]+" "+act[0]+" "+strconv.FormatBool(res))
	}else{
		fmt.Fprintf(w,"验证失败")
		fmt.Fprintf(w,		strconv.FormatBool(res))
	}
}


