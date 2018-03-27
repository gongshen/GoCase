* 只有当defer语句执行完，外围函数才会返回
* 外围函数如果引发了`panic`，`defer`语句也会执行完panic才会扩散
* 在defer执行的时候，针对defer语句的表达式会被压栈，等到外围函数结束时，才依次从栈中取出
* defer是在函数退出前执行，并且是`值传递`
# 1、闭包的值传递和指针传递
### 1.闭包的值传递
```go
func main()  {
	var m int
	m=10
	fmt.Println(m)
	defer func(i int) {
		fmt.Println(i)
	}(m)
	m=20
}
```
结果是：
```go
10
10
```
### 2.闭包的指针传递:
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

# 2、defer的值传递
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
-----
# 3、容易入的"坑"
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