package main

import "fmt"

func SelectSort(a *[5]int)  {
	for j:=0;j<len(a)-1 ;j++ {
		max := a[j]
		maxIndex := j
		for i := j + 1; i < len(a); i++ {
			if max < a[i] {
				max = a[i]
				maxIndex = i
			}
			a[j], a[maxIndex] = a[maxIndex], a[j]
		}
	}
}
func main()  {
	var a [5]int=[5]int{55,66,33,44,99}
	SelectSort(&a)
	fmt.Print(a)
}
