package main
//wait group
import (
	"net/http"
	_"runtime"
	_"time"
	"sync"
)
import "fmt"

func print11()  {
	for i:=0;i<1000;i++{
		fmt.Print(11," ")
	}

}
func print7()  {
	for i:=0;i<1000;i++{
		fmt.Print(7," ")
	}
}
func print3()  {
	for i:=0;i<1000;i++{
		fmt.Print(3," ")
	}
}
func consumer(c chan int)  {

	
}
func producter(c chan int)  {
	sum:=0
	for i:=0;i<10000;i++{
		sum+=i
	}
	fmt.Print(sum)
	c<-sum

}
func getHTTP(url string)  {
	res,err:=http.Get(url)
	if(err!=nil){
		fmt.Print("error")
	}
	fmt.Print(url," ",res.Status,"\n")
	defer wg.Done()

}
var wg sync.WaitGroup
var urls= []string{"http://www.baidu.com","http://www.qq.com","http://sogou.com"}
func main(){

	for _,v:=range urls{
		wg.Add(1)
		go getHTTP(v)
	}
	wg.Wait()


}
