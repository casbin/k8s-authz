package main

import "fmt"
import "time"

func fibonaqi( n int) int   {
	if n<2 {
		return n
	}else {
		 a:= fibonaqi(n-1)+fibonaqi(n-2)

		 return a
	}
	
}
func spin(duration time.Duration){
	for{
		for _,r:=range `/|\-` {
			fmt.Printf("\r%c", r)
			time.Sleep(duration)
		}
	}


}
func main()  {
	go spin(100*time.Millisecond)
	 fmt.Print("\n",fibonaqi(40))



}