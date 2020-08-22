package main

import (
	"fmt"
	"math/rand"
	"runtime"
)

func product(done chan struct{}) chan int {
	ch:=make(chan int)

	go func() {
		Label:for{
				select {
					case ch<-rand.Int():

					case <-done:
						break Label

				}

		}
		close(ch)
	}()
	return ch

}
func main()  {
	done:=make(chan struct{})
	ch:=product(done)

	fmt.Print(<-ch,"\n")

	fmt.Print(<-ch,"\n")

	fmt.Print(<-ch,"\n")

	close(done)
	fmt.Print(<-ch,"\n")

	fmt.Print(<-ch,"\n")

	fmt.Print(<-ch,"\n")
	fmt.Print(runtime.GOMAXPROCS(0))

}
