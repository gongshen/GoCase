### 不定长变参其实是一个切片，可以range遍历
```go
// ...interface{} 表示将参数形成一个切片
func TestArgs(first int,arg ...interface{}){
    fmt.Println(first,arg)
}
func main(){
    nums:=[]int{1,2,3,4,5,6}
    //TestArgs(1,nums)
    TestArgs(1,nums...)                //表示将切片打散再传入
}
//cannot use nums (type []int64) as type []interface {} in argument to TestArgs
```
以上代码出现了类型不匹配的错误
**原因：**
因为是直接将slice传进去，类型不匹配。
```go
func TestArgs(first int,arg ...interface{}){
    fmt.Println(first,arg)
}
func main(){
    nums:=[]interface{}{1,2,3,4,5}
    TestArgs(1,nums)        
}
//1 [[1,2,3,4,5]]
```
**小结：**
- TestArgs(1,nums...)    //将nums打散再传入
- 使用 ...语法糖的slice时，直接传入这个slice
- 单个可变参数实际是执行[]T{arg1,arg2}