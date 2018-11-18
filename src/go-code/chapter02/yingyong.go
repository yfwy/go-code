package main
import "fmt"
func main(){
	var arr [26] byte
	for i:=0;i<len(arr);i++{
		arr[i] =  'A'+ byte(i)
	}
	fmt.Printf("%c",arr)
}