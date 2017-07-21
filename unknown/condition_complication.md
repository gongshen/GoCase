在源文件(.go, .h, .c, .s 等)头部添加 "+build" 注释，指示编译器检查相关环境变量。多个约束标记会合并处理。其中空格表示 OR，逗号 AND，感叹号 NOT。
```go
// +build windows linux 						<-- (表示支持windows或者linux编译)AND(amd64和不用cgo)
// +build amd64,!cgo							
												<-- 必须有空行，区别包文档
package main
```