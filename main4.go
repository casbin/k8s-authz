package main

import (
	"fmt"
	"math/rand"
)
//带缓冲, 并发 ,退出通知机制的生成器
func GenerateIntA(done chan struct{}) chan int {
	ch:=make(chan int,10)
	go func() {
		Label:
		for{
			select {
			case ch<- rand.Int():
			case <-done:
				break Label

			}
		}
		close(ch)

	}()
	return ch
}
func GenerateIntB(done chan struct{}) chan int {
	ch:=make(chan int,10)
	go func() {
	Label:
		for{
			select {
			case ch<- rand.Int():
			case <-done:
				break Label

			}
		}
		close(ch)

	}()
	return ch
}

func GenerateInt(done chan struct{})chan int {
	ch:=make(chan  int ,20)
	notify:=make(chan struct{})
	go func() {
		Label:
			for {
				select {
				case ch <- <-GenerateIntA(notify):
				case ch <- <-GenerateIntB(notify):
				case <-done:
					notify <- struct{}{}
					notify <- struct{}{}
					break Label


				}
			}
		close(ch)
	}()
	
	return  ch
}
func main()  {
	done:=make(chan struct{})

	ch:=GenerateInt(done)

	for i:=1;i<20;i++{
		fmt.Print(<-ch,"\n")
	}
	done<- struct{}{}

	fmt.Print("stop generate")

}