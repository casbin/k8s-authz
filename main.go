package main

import (
	"bufio"
	"fmt"
	"os"
	"io/ioutil"
	"strings"
)

func main() {
	counts:=make(map[string]int)

	for _, filename :=range os.Args[1:] {
		data, err := ioutil.ReadFile(filename)

		if err!=nil{
			fmt.Println("Read file err")
		}
		for _,line := range strings.Split(string(data),"\n"){
			counts[line]++
		}
		for key,value :=range counts{
			if value>1 {
				fmt.Println(filename)
				delete(counts, key)
				break
			}
		}


	}


}

//
func ReadFromFile(){
	counts  :=make(map[string]int)
	files:=os.Args[1:]
	if len(files)==0{
		countLines(os.Stdin,counts)
	}else{
		for _,arg:= range files{
			f,err :=os.Open(arg)
			if err!=nil{
				fmt.Println("打开文件错误")

			}

			countLines(f,counts)
			f.Close()
		}

		for  line, _:=range counts{
			fmt.Println(line)
		}


	}

}
//函数和包级别的变量和声明, 可以任意顺序声明
func countLines(f * os.File,counts map[string]int)  {
	input :=bufio.NewScanner(f)
	for input.Scan(){
		counts[input.Text()]++
	}

}