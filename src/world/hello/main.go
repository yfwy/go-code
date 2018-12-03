package main

import (
	"github.com/astaxie/beego"
	_ "hello/models"
	_ "hello/routers"
)

func main() {
	beego.AddFuncMap("showprepage",prepage)
	beego.AddFuncMap("shownextpage",shownextpage)
	beego.Run()
}
/*func HandleParper(data int) string {
	paperIndex := data-1
	paperIndex1 :=strconv.Itoa(paperIndex)
	return paperIndex1
}*/
func prepage(pageindex int)(preIndex int){
	preIndex = pageindex - 1
	return
}

func shownextpage(pageindex int)(nextIndex int){
	nextIndex = pageindex + 1
	return
}
