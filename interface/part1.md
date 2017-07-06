myInterface接口有一个Print方法，myStruct结构实现了Print方法，也就等于实现了myInterface接口。
```go
package main

import "fmt"

type myInterface interface {
	Print()
}

func TestFunc(x myInterface)  {
	fmt.Println("Ok!")
}

type myStruct struct {}

func (me *myStruct)Print(){}

func main()  {
	var mi *myStruct
	TestFunc(mi)
}
```