package main

import (
	"fmt"
	"strconv"
)


type account struct {
	name int
	mima int
	money int
}

func Newaccount(name int,mima int,money int) *account {
	if len(strconv.Itoa(name))>10||len(strconv.Itoa(name))<6{
	fmt.Print("accont false")
	 return nil
	}
	if mima!=6 {
		fmt.Print("mama   false")
		return nil
	}
	if money<20 {
		fmt.Print("money false")
		return nil
	}
	return &account{
		name:name,
		mima:mima,
		money:money,
	}

	}