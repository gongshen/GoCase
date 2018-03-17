package main

import (
	"os"
	"fmt"
)


func main() {
	var s,sep string
	for _,arg:=range os.Args[1:]{
		s+=sep+arg
	}
	fmt.Println(s)
	/*
	for i:=0;i<len(os.Args);i++{
		s+=sep+os.Args[i]
		sep=" "
	}
	*/
}

