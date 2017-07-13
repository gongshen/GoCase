# 1.go的内联优化
```go
package main

func Max(a,b int)int {
	if a>b {
		return a
	}else{
		return b
	}
}

func DoubleMax(a,b int) int {
	return 2*Max(a,b)
}

func main()  {
	DoubleMax(1,2)
}
```
使用内联优化
```shell
$ go build -gcflags "-m" -o test main.go
# command-line-arguments
.\main.go:5: can inline Max
.\main.go:13: can inline DoubleMax
.\main.go:14: inlining call to Max
.\main.go:17: can inline main
.\main.go:18: inlining call to DoubleMax
.\main.go:18: inlining call to Max
```
# 2.条件编译
在源文件(.go, .h, .c, .s 等)头部添加 "+build" 注释，指示编译器检查相关环境变量。多个约束标记会合并处理。其中空格表示 OR，逗号 AND，感叹号 NOT。
```go
// +build windows linux 						<-- (表示支持windows或者linux编译)AND(amd64和不用cgo)
// +build amd64,!cgo							
												<-- 必须有空行，区别包文档
package main
```
# 3.数据竞争检查
```go
func main()  {
	var wg sync.WaitGroup
	wg.Add(2)
	x:=100
	go func() {
		defer wg.Done()
		for	{
			x+=1
		}
	}()
	go func() {
		defer wg.Done()
		for {
			x+=100
		}
	}()
	wg.Wait()
}
```
数据竞争严重影响性能
```shell
$ GOMAXPROCS=2 go run -race main.go
==================
WARNING: DATA RACE
Read at 0x00c042038010 by goroutine 6:
  main.main.func2()
      D:/gopath/src/main.go:18 +0x6c

Previous write at 0x00c042038010 by goroutine 5:
  main.main.func1()
      D:/gopath/src/main.go:12 +0x85

Goroutine 6 (running) created at:
  main.main()
      D:/gopath/src/main.go:20 +0x10f

Goroutine 5 (running) created at:
  main.main()
      D:/gopath/src/main.go:14 +0xe3
==================
```