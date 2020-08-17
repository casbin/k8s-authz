package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
 	"os"
	"sync"
)

var  mu sync.Mutex
var	count int
func  main()  {
	http.HandleFunc("/",handler)
	http.HandleFunc("/count",counter)

	log.Fatal(http.ListenAndServe("localhost:8000",nil))
}

func handler(w http.ResponseWriter,r *http.Request)  {

	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w,"URL path access by you is %q", r.URL.Path)

}
func counter(w http.ResponseWriter,r* http.Request)  {
	mu.Lock()
	fmt.Fprintf(w,"現在有%v個人訪問過此頁面",count)
	mu.Unlock()

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