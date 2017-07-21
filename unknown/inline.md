```go
package main

func test()*int  {
	x:=new(int)
	*x=0xAABB
	return x
}

func main()  {
	println(test())
}
```
禁用内联优化:两个栈帧间需要传递对象，所以在`堆`上分配。
```shell
$ go build -gcflags "-l" -o test main.go
$ go tool objdump -s "main\.test" test
...
main.go:4       0x44dedf        e80ce4fbff              CALL runtime.newobject(SB)
...
```
使用内联优化后:test函数就等于是`main栈帧`内的局部变量，无需`堆`分配，从而提升性能。
```shell
$ go build -o test main.go
$ go tool objdump -s "main\.main" test 
```
