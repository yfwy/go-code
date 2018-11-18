package main

import (
"fmt"
"github.com/garyburd/redigo/redis"
)
func main()  {
	conn,err:=redis.Dial("tcp" ,"127.0.0.1:6379")
	if err != nil {
		fmt.Print("dds")
		return
	}
	_, err = conn.Do("HSet", "name", "tom","555")
	if err!=nil {
		fmt.Printf("cuowu")
	}

	fmt.Println(conn)
}