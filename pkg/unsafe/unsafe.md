## 1.从slice中获取一块内存空间
```go
s:=make([]int,10)
v:=unsafe.Pointer(&s[0])
```
## 2.从内存构造slice
```go
var ptr unsafe.Pointer
s:=(*[1<<20]byte)(unsafe.Pointer(ptr))[:200]
//构建了一个200元素的slice
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
## 6.测试[]byte普通转化string和伪造string的性能差异：
```go
package main

import (
	"testing"
	"unsafe"
)

//测试[]byte转化为string
func Test_ByteString(t *testing.T)  {
	var x=[]byte("Hello World!")
	var y=*(*string)(unsafe.Pointer(&x))
	var z=string(x)
	if z!=y{
		t.Fail()
	}
}

func Benchmark_Normal(b *testing.B)  {
	var x=[]byte("Hello World!")
	for i:=0;i<b.N;i++{
		_=string(x)
	}
}

func Benchmark_ByteString(b *testing.B)  {
	var x=[]byte("Hello World!")
	for i:=0;i<b.N;i++{
		_=*(*string)(unsafe.Pointer(&x))
	}
}
```
结果对比：
因为普通转化需要拷贝
```shell
$ go test -bench .
Benchmark_Normal-4              100000000               10.4 ns/op
Benchmark_ByteString-4          2000000000               0.37 ns/op
PASS
ok      Exp3  1.868s
```
## 7.unsafe.Offsetof()和unsafe.Sizeof()
`Offsetof`参数字段必须是`x.f`,返回f相对于x的偏移量。
`Sizeof`表示操作数在内存中的字节大小。
```go
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var s= struct {
		a int
		b int
		c int
		d int
	}{0,0,0,0}
	p:=unsafe.Pointer(&s)
	pb:=(*int)(unsafe.Pointer(uintptr(p)+unsafe.Offsetof(s.c)))
	*pb=20
	fmt.Println(s)
}
```
结果：
```shell
{0 0 20 0}
```
## 8.一个struct变量的内存地址分配
一个结构体变量 x 以及其在64位机器上的典型的内存. 灰色区域是空洞.
```go
var x struct{
	a bool
	b int16
	c []int
}
```
![struct](https://github.com/gongshen/GoCase/blob/master/pic/struct.png)
```go
package main

import (
	"fmt"
	"unsafe"
)

func main() {

	var x struct {
		a bool
		b int16
		c []int
	}
	fmt.Printf("%-30s%-30s%-30s%-50s\n",
		"Row", "Sizeof", "Alignof(对齐倍数)", "Offsetof(偏移量)")
	fmt.Printf("%-30s%-30d%-30d%-50s\n",
		"x", unsafe.Sizeof(x), unsafe.Alignof(x), "")
	fmt.Printf("%-30s%-30d%-30d%-50d\n",
		"x.a", unsafe.Sizeof(x.a), unsafe.Alignof(x.a), unsafe.Offsetof(x.a))
	fmt.Printf("%-30s%-30d%-30d%-50d\n",
		"x.b", unsafe.Sizeof(x.b), unsafe.Alignof(x.b), unsafe.Offsetof(x.b))
	fmt.Printf("%-30s%-30d%-30d%-50d\n",
		"x.c", unsafe.Sizeof(x.c), unsafe.Alignof(x.c), unsafe.Offsetof(x.c))
}
```
结果：
```go
Row                           Sizeof                        Alignof(对齐倍数)                 Offsetof(偏移量)                                     
x                             32                            8                                                                               
x.a                           1                             1                             0                                                 
x.b                           2                             2                             2                                                 
x.c                           24                            8                             8                                                 

Process finished with exit code 0
```