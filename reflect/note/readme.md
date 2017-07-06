1. 如果要调用Elem()方法，那么必须是`接口`或者`指针`。
2. `a...`表示将参数切片a打散再传入
```go
func (u User) Hello(str string,a ...interface{})string{
	return fmt.Sprintf(str,a...)
}
```
