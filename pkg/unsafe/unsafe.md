## 1.从slice中获取一块内存空间
```go
s:=make([]int,10)
v:=unsafe.Pointer(&s[0])
```