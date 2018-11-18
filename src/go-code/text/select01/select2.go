package main

import "fmt"

func SelectSort(a *[5]int)  {
	for j:=1;j<len(a) ;j++ {
		insertVal := a[j]
		insertIndex := j - 1
		for insertIndex >= 0 && a[insertIndex] < insertVal {
			a[insertIndex+1] = a[insertIndex]
			insertIndex--
		}
		if insertIndex+1 != j {
			a[insertIndex+1] = insertVal
		}
	}
}

func main()  {
	var a [5]int=[5]int{55,66,33,44,99}
	SelectSort(&a)
	fmt.Print(a)
}
