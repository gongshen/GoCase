## defer是在函数退出前执行，并且是值传递。
```go
package main

import "fmt"

func main()  {
	var m int
	m=10
	fmt.Println(m)
	defer Print(m)
	m=20
}
func Print(m interface{})  {
	fmt.Println(m.(int))
}
```
结果是：
```go
10
10
```
如果用指针作为函数参数：
```go
package main

import "fmt"

func main()  {
	var m *int
	m=new(int)
	*m=10
	fmt.Println(*m)
	defer Print(m)
	*m=20
}
func Print(m interface{})  {
	fmt.Println(*m.(*int))
}
```
结果是：
```go
10
20
```
也可以使用`闭包`:
```go
package main

import (
	"fmt"
)

func main()  {
	var m int
	m=10
	fmt.Println(m)
	defer func() {
		fmt.Println(m)				//闭包的意思就是他已经把指针准备好了
	}()
	m=20
}
```
结果是：
```go
10
20
```