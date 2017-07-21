# 1、下面的代码有没有什么问题？
```go
package main

import (
	"fmt"
)

type student struct {
	Name string
	Age  int
}

func pase_student() map[string]*student {
	m := make(map[string]*student)
	stus := []student{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}
	for _,stu:=range stus{
			m[stu.Name]=&stu
	}
	return m
}
func main() {
	students := pase_student()
	for k, v := range students {
		fmt.Printf("key=%s,value=%v \n", k, v)
	}
}
```
**结果:**
```shell
key=li,value=&{wang 22} 
key=wang,value=&{wang 22} 
key=zhou,value=&{wang 22} 
```
**解析:**
```go
for _,stu:=range stus{
			m[stu.Name]=&stu
	}
```
> 因为for range遍历的时候，stu是值拷贝，stu变量的指针不变，所以取stu的地址都是同一个，很显然最后遍历的wang 22就是stu的值。
**修正方案一：**
取数组中原始的指针：
```go
for i,_:=range stus{
	stu:=stus[i]
	m[stu.Name]=&stu
}
```
**修正方案二：**
不要取stu的地址，直接进行值拷贝，当然map的键值对也要改下。
```go
func pase_student() map[string]student {
	m := make(map[string]student)
	stus := []student{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}
	for _,stu:=range stus{
			m[stu.Name]=stu
	}
	return m
}
```