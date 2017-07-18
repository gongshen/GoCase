# Defer和闭包
> defer是在函数退出前执行，并且是`值传递`。
### 1.函数参数是值传递：
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
### 2.如果用指针作为函数参数：
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
### 3.也可以使用`闭包`:
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
### 例子1：
```go
func f()(result int){
	defer func(){
		result++
	}()
	return 0
}
```
> 解：
`return 0`并不是原子操作,应该写成`result=0`,`return`;所以结果是`1`。

### 例子2：
```go
func f()(r int){
	defer func(r int){
		r=r+5
	}(r)
	return 1
}
```
> 解：
先`r=1`,再`return`,但是闭包中是值传递,不会影响r的值,所以结果是`1`。