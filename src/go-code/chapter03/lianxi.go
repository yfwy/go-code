package main

import (
	"fmt"
)

func gongNeng(arr map[string]map[string]string,name string)  {
		if arr[name]!=nil {
		arr[name]["mima"]="88888888"
		}else {
			arr[name]=make(map[string]string)
			arr[name]["mima"]="88888888"
			arr[name]["nicheng"]="a go"+name
		}
		fmt.Print(arr)
	}
func main(){
	arr :=make(map[string]map[string]string)
	gongNeng(arr,"marine")
}