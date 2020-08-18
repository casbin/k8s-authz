package main

import (
	"fmt"
	"log"
	"net/http"
)
//笔记
//listenAndServe这个函数第一个参数是一个url,第二个参数是一个 handler的借口
/*
/所以,我们要写一个类,实现handler这个借口
接口里面有一个方法,叫ServeHTTP 实现了这个函数,就实现了这个接口
 */
type dollor float64
type db map[string]dollor

func (d db)list (w http.ResponseWriter,r *http.Request){
	for k,v :=range d{
				fmt.Fprintf(w,"%v %v\n",k,v)
	}
}

func (d db)price (w http.ResponseWriter,r * http.Request){
	for k,v :=range d{
		fmt.Fprintf(w,"%v %v\n",k,v)
	}
	for k,v :=range d{
		fmt.Fprintf(w,"%v %v\n",k,v)
	}
	for k,v :=range d{
		fmt.Fprintf(w,"%v %v\n",k,v)
	}
}
func main(){
	d:=db{"iphone 11":4900,"ipad 2019":2019}
	mux:=http.NewServeMux()
	//mux.Handle("/price",http.HandlerFunc(d.price))
	mux.HandleFunc("/price",d.list)
	mux.HandleFunc("/list",d.price)
	//mux.Handle("/list",http.HandlerFunc(d.list))
	log.Fatal(http.ListenAndServe("localhost:8000",mux))


}