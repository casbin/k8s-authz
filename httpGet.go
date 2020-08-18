package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
 	"os"
	"sync"
	"flag"
)

var  mu sync.Mutex
var	count int
var sep=flag.String("s"," ","-s print the name of you")
var n = flag.Bool("n",false,"-n stdout<<shuchu1")
func  main()  {
	flag.Parse()
	if *n!=false {
		fmt.Println("shuchu1")
	}
	http.HandleFunc("/",handler)
	http.HandleFunc("/count",counter)

	log.Fatal(http.ListenAndServe("0.0.0.0:8000",nil))
}

func handler(w http.ResponseWriter,r *http.Request)  {

	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w,"URL path access by you is %q", r.URL.Path)

}
func counter(w http.ResponseWriter,r* http.Request)  {
	mu.Lock()
	fmt.Fprintf(w,"現在有%v個人訪問過此頁面\n",count)
	mu.Unlock()

	fmt.Fprintf(w,"%s %s %s\n",r.Method,r.URL,r.Proto)

	for k,v:=range r.Header{
		fmt.Fprintf(w,"Header[%q] =%q\n",k,v)

	}

	fmt.Fprintf(w,"Host=%q\n",r.Host)
	fmt.Fprintf(w,"RemoteAddr=%q\n",r.RemoteAddr)

	if err:=r.ParseForm();err!=nil{
		log.Println(err)
	}

	for k,v :=range r.Form{
		fmt.Fprintf(w,"Form[%q]=%q\n",k,v)
	}




}














































func network(){


	for _,url := range os.Args[1:] {
		if url[0:7]!="http://" {
			url="http://"+url
		}
		res,err:=http.Get(url)
		b:=res.Status
		if err!=nil {
			fmt.Fprintf(os.Stderr,"fetch: %v\n http状态码等于%s",err,b)
			os.Exit(1)
		}

		io.Copy(os.Stdout,res.Body)
		res.Body.Close()
		if  err!=nil{
			fmt.Fprintf(os.Stderr,"fetch  :reading %s %v",url,err)
			os.Exit(1)
		}

	}
}