# 可寻址值和不可寻址值调用方法
对于receiver是指针类型的方法，可寻址值的变量能直接调用该方法。
那么对于不可寻址值得变量自然不能直接调用。
```go
type data struct {
	name string
}

func (d *data)Print()  {
	fmt.Println("name",d.name)
}

type Printer interface {
	Print()
}
```
```go
func main{
	d1 := &data{"one"}
	d1.Print()
}
```
> `d1`是可寻址变量，能直接调用`Print`方法。

```go
func main{
	var p Printer=&data{"two"}		//因为*data实现了该方法，所以要取址
	p.Print()
}
```
> `p`是可寻址变量，能直接调用`Print`方法。

```go
func main{
	m1:=map[string]data{"x":{"three"}}
	n:=m1["x"]
	n.Print()
}
```
> `map`中的`Value`是一个不可寻址的变量，所以不能直接调用`Print`

```go
func main{
	m2:=map[string]data{"x":{"four"}}
	m2.Print()
}
```
> 因为`Value`是指针类型的`*data`，所以能直接调用`Print`。

```go
func main{
	s:=[]data{
		{"five"},
	}
	s.Print()
}
```
> 对于结构体类型的`slice`是可以直接调用的。
---
请看下面完整代码：
```go
package main

import "fmt"

type data struct {
	name string
}

func (d *data)Print()  {
	fmt.Println("name",d.name)
}

type Printer interface {
	Print()
}

func main() {
	d1 := &data{"one"}
	d1.Print()

	var p Printer = &data{"two"}
	p.Print()
	
	m1 := map[string]data{"x":{"three"}}
	n:=m1["x"]
	n.Print()

	m2:=map[string]*data{"x":{"four"}}
	m2["x"].Print()

	s:=[]data{
		{"five"},
	}
	s[0].Print()
}
```