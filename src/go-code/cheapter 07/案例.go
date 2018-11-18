package cheapter_0

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Monster struct {
	Name string
	Age int
	Skill string
}

func (this *Monster)Store() {
	 Mon :=Monster{
	 	Name:"牛魔王",
	 	Age:550,
	 	Skill:"红孩儿",
	 }
	 data,err:=json.Marshal(Mon)
	if err!=nil {
		fmt.Print("you")
	}
	 filePath:= "d:/monster.ser"
	 err=ioutil.WriteFile(filePath,data,06666)
	if err!=nil {
		fmt.Print("false")
	}
}
func (this *Monster)ReStore()  {
	filePath:="d:/monster.ser"
	data,err:=ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Print("youcuo" )
	}
	err =json.Unmarshal(data,this)
	if err != nil {
		fmt.Print("youwu2")
	}

	}