package main

import (
	"fmt"
	"io/ioutil"
)

func main(){
	file:="c:/dd/text01.tet"
	c,err:=ioutil.ReadFile(file)
	if err != nil {
		fmt.Print("dddd")
	}
	fmt.Printf("%v",c)//qie   pian[]byte
}