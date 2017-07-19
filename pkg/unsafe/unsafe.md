## 1.从slice中获取一块内存空间
```go
s:=make([]int,10)
v:=unsafe.Pointer(&s[0])
```
## 2.从内存构造slice
```go
var ptr unsafe.Pointer
s:=(*[1<<20]byte)(ptr)
```
## 3.使用reflect.SliceHeader来构造slice
```go
var ptr []byte
	x:=(*reflect.SliceHeader)(unsafe.Pointer(&ptr))
	x.Len = 10
	x.Cap = 10
	x.Data= uintptr(unsafe.Pointer(&ptr))
	res:=*(*[]byte)(unsafe.Pointer(&x))
```
结构体与[]byte的互转：
```go
func main() {
	s:=&MyStruct{100,200}
	fmt.Printf("%v\n",MyStructToBytes(s))
}

type MyStruct struct {
	A int
	B int
}

func MyStructToBytes(s *MyStruct) []byte {
	var x reflect.SliceHeader
	x.Len = 10
	x.Cap = 10
	x.Data = uintptr(unsafe.Pointer(s))
	return *(*[]byte)(unsafe.Pointer(&x))
}
```
## 4.实现slice的底层结构，再构造slice
```go
var ptr unsafe.Pointer
	var s = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(ptr), 10, 10}
	slice := *(*[]byte)(unsafe.Pointer(&s))
	fmt.Printf("%T\n",slice)
```
## 5.导出未导出的变量
lib/lib.go
```go
package lib

type Student struct {
	x int
	Y int
}

func NewStu() *Student {
	return new(Student)
}
```
main.go
```go
package main

import (
	"./lib"
	"unsafe"
	"fmt"
)

func main()  {
	s:=lib.NewStu()
	s.Y = 100
	p:=(*struct{x int})(unsafe.Pointer(s))
	p.x= 200
	fmt.Println(s)
}
```