package main

import (
	"fmt"
)

type HeroNode struct {
	name string
	no int
	next *HeroNode
}

func Add(head *HeroNode, NewHero *HeroNode)  {
	temp:=head
	for{
		if temp.next==nil {
			break
		}
		temp = temp.next
	}
	temp.next =NewHero
}
func List(head *HeroNode){
	temp:=head
	if temp.next == nil {
		fmt.Print()
	}
	for {
		if temp.next == nil {
			return
		}
		fmt.Print(temp.no)
		temp=temp.next
	}
}
func InsertHero(head *HeroNode,newHero *HeroNode)  {
	temp:=head
	for{
		if temp.next == nil {
			break
		}else if newHero.no>temp.next.no {
			break
		}else if newHero.no==temp.next.no{
			return
		}
		newHero.next=temp.next
		temp.next=newHero
	}
}
func main() {
	Head :=&HeroNode{}
	hero01 := &HeroNode{
		name:"songjiang",
		no:1,
	}
	hero02 :=&HeroNode{
		name:"lu jvyi",
		no:2,
	}
	hero03 :=&HeroNode{
		name: "jjj",
		no:   3,
	}
}