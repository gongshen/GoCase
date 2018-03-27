* panic用于停止当前的控制流程并引发一个恐慌
* recover用于使当前程序从恐慌中恢复并重新获得流程控制权
* recover函数的结果是一个interface{}类型，如果结果不是nil，那么就有问题啦！
* recover和defer应该配合使用

我们看下标准库fmt中的Token是怎么处理的
```go
func (s *ss) Token(skipSpace bool, f func(rune) bool) (tok []byte, err error) {
	defer func() {
		if e := recover(); e != nil {		//判断recover函数的结果是否为nil
			if se, ok := e.(scanError); ok {	//判断panic的类型
				err = se.err		
//如果panic是这个类型，那么这个值就会赋值给结果值变量err，
//这样做到了精确控制panic，将已经recover的恐慌当作常规结果返回
			} else {
				panic(e)
//否则恐慌会再次引发
			}
		}
	}()
...
}
```
恐慌被传递到调用栈的最顶层的结果：
```go
panic: An intended fatal error! [recovered]
	panic: An intended fatal error!
```
