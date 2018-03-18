package main

import (
	"time"
	"fmt"
)

func calculate()  {
	t1:=time.Now()	// get current time
	// start calculate
	for i:=0;i<1000;i++{
		fmt.Print("*")
	}
	elapsed:=time.Since(t1)
	fmt.Println("\nApp elapsed:",elapsed)
}
func main(){
	calculate()
}