package main

import (
	"fmt"
)

type ValNode struct {
	row int
	col int
	val int
}

func main() {
	/*var a [11][11] int
	a[2][3] = 1
	a[3][2] = 2
	for _ ,v:=range a{
		for _ ,v:=range v{
			fmt.Print(v)
		}
		fmt.Println()
	}
	fmt.Print()*/
	var a [11][11] int
	a[2][3]=2
	a[3][3]=3
	for i, v:=range a {
		for j, v2 := range v {
			if v2 != 0 {
				var val = ValNode{
					row: i,
					col: j,
					val: v2,
				}
				fmt.Print(val)
			}
		}
	}

}