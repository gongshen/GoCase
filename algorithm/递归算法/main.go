package main

import (
	"fmt"
)

func fn(i int)int{
	if i<=1{
		return 1
	}else{
		return i*(fn(i-1))
	}
}

func main()  {
	fmt.Println(fn(8))
}