package itface

import "fmt"
import _"encoding/json"
type Interface interface{
	Smile() int
}
type Person struct{
	Name string`json:"name"`
	Age int `json:"age"`
	
}


func (p Person) Smile() int {
	fmt.Println(p.Name+ " smile")
	return 1
}
func Do(i interface{}){
	fmt.Print(i)
}