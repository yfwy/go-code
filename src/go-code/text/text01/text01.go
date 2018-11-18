package main

import (
	"fmt"
)

type Heronode struct {
	name string
	no int
	pre *Heronode
	next *Heronode
}

func InsertHero(head *Heronode,newHero *Heronode)  {
	temp:=head
	for{
		if temp.next == nil {
			break
		}
		temp= temp.next
	}
	temp.next=newHero
	newHero.pre=temp
}
func listHero(head *Heronode) {
	temp:=head
	if temp.next == nil {
		fmt.Print("cuowu" )
		return
	}
	for{
		fmt.Println(temp.next.name)
		temp=temp.next
		if temp.next == nil {
			break
		}
		
	}
}
func  InsertHero2(head*Heronode,hero*Heronode){
	temp:=head
	for{
		if temp.next == nil {
			break
		}else if hero.no<temp.next.no{
			break
		}else if hero.no==temp.next.no {
			fmt.Scanln("cuowu" )
			return
		}
		temp=temp.next
	}
	hero.next=temp.next
	temp.next=hero
}
func main(){
	Head :=&Heronode{}
	hero01 := &Heronode{
		name:"songjiang",
		no:1,
	}
	hero02 :=&Heronode{
		name:"lu jvyi",
		no:2,
	}
	hero03 :=&Heronode{
		name:"jjj",
		no:3,
	}
	InsertHero(Head,hero01)
	InsertHero2(Head,hero03)
	InsertHero2(Head,hero02)

	listHero(Head)

}