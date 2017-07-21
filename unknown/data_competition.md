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