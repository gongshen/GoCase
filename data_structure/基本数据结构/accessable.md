# 可寻址值和不可寻址值调用方法
### 示例一：
对于一个非指针类型，它关联的方法集中只包含它的值方法。对于一个指针类型，它关联的方法集中既包含值方法和指针方法。
但是非指针类型也是可以调到指针方法的，因为Go进行了自动转换
---

```go
type data struct {
	name string
}
// 指针类型的方法
func (d *data)Print()  {
	fmt.Println("name:",d.name)
}
// 值类型的方法
func (d data)Print2(){
	fmt.Println("name2:",d.name)
}

type Printer interface {
	Print()
}

func main() {
	var p Printer = &data{"one"}
	p.Print()

// p是一个指针类型

	m1 := map[string]data{"x":{"two"}}
	n:=m1["x"]
	n.Print()

// 重要！
// n是值类型，但是可以调用指针类型的方法，Go进行了自动转换
// n.Print()  ==  (&n).Print()

	m2:=map[string]*data{"x":{"three"}}
	n2:=m2["x"]
	n2.Print2()
	
// 重要！
// n2是指针类型，但是可以调用值类型的方法，Go进行了自动转换
// n2.Print2()  ==  (*n2).Print2()

	s:=[]data{
		{"four"},
	}
	s[0].Print()
	
// 对于结构体类型的`slice`是指针类型
}
```
### 示例二：
receiver变量其实就是源值的一个复制品。如果receiver是值类型，那么自然没有办法修改源值；如果receiver是指针类型，那么指针值指向的就是源值的地址，就能够修改源值。
---

```go
type myInt int

// 指针类型方法
func (i *myInt)add2(another int)myInt{
	*i=*i+myInt(another)
	return *i
}
// 值类型方法
func (i myInt)add(another int)myInt{
	i=i+myInt(another)
	return i
}

func main(){
	i1:=myInt(1)
	i2:=i1.add(2)
	fmt.Println(i1,i2)
// 结果是：1,3
	i2=i1.add2(2)
	fmt.Println(i1,i2)
// 结果是：3,3
}
```