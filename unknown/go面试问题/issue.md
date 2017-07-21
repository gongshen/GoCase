# 1、下面的代码有没有什么问题？
```go
package main

import (
	"fmt"
)

type student struct {
	Name string
	Age  int
}

func pase_student() map[string]*student {
	m := make(map[string]*student)
	stus := []student{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}
	for _,stu:=range stus{
			m[stu.Name]=&stu
	}
	return m
}
func main() {
	students := pase_student()
	for k, v := range students {
		fmt.Printf("key=%s,value=%v \n", k, v)
	}
}
```
**结果:**
```shell
key=li,value=&{wang 22} 
key=wang,value=&{wang 22} 
key=zhou,value=&{wang 22} 
```
**解析:**
```go
for _,stu:=range stus{
			m[stu.Name]=&stu
	}
```
> 因为for range遍历的时候，stu是值拷贝，stu变量的指针不变，所以取stu的地址都是同一个，很显然最后遍历的wang 22就是stu的值。

**修正方案一：**
取数组中原始的指针：
```go
for i,_:=range stus{
	stu:=stus[i]
	m[stu.Name]=&stu
}
```
**修正方案二：**
不要取stu的地址，直接进行值拷贝，当然map的键值对也要改下。
```go
func pase_student() map[string]student {
	m := make(map[string]student)
	stus := []student{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}
	for _,stu:=range stus{
			m[stu.Name]=stu
	}
	return m
}
```
# 2、下面代码有什么输出，为什么。
```go
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("i: ", i)
			wg.Done()
		}()
	}
	
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("i: ", i)
			wg.Done()
		}(i)
	}
	
	wg.Wait()
}
```
**解析:**
> 第一个for循环遍历：因为闭包默认帮你实现了指针，所以遍历结束，i=10.
第二个for循环遍历中：因为闭包传入了参数，实现了值拷贝，所以i=0，1，2，... ，9。

# 3、下面的代码会触发异常吗？
```go
func main() {
	runtime.GOMAXPROCS(1)
	int_chan := make(chan int, 1)
	string_chan := make(chan string, 1)
	int_chan <- 1
	string_chan <- "hello"
	select {
	case value := <-int_chan:
		fmt.Println(value)
	case value := <-string_chan:
		panic(value)
	}
}
```
**解析:**
> 可能会发生panic。
因为channel引入了缓存，当select下的case有多个时，随机执行其中一个，所以channel能接受数据也能发生数据。

# 4、下面代码的输出：
```go
package main

import (
	"fmt"
)

func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

func main() {
	a := 1                                             	//line 1
	b := 2                                             	//2
	defer calc("1", a, calc("10", a, b))  		 //3
	a = 0                                              	//4
	defer calc("2", a, calc("20", a, b))  		 //5
	b = 1                                              	//6
}
```
**输出:**
```shell
10 1 2 3
20 0 2 2
2 0 2 2
1 1 3 4
```
**解析:**
> 1.defer语句是return原子操作之前执行，所以`1`和`2`最后执行。
2.因为defer是值传递，所以第3行中的`a`和`b`分别是`1`和`2`；第5行中的`a`和`b`分别是`0`和`2`。

# 5、下面代码的输出：
```go
package main

import (
	"fmt"
)
func main() {
	s := make([]int, 5)
	s = append(s, 1, 2, 3)
	fmt.Println(s)
}
```
**输出:**
```shell
[0 0 0 0 0 1 2 3]
```
**解析:**
> 很多人以为是`[1 2 3]`，但是当你执行`make([]int,5)`的时候，已经将5个int类型的数据初始化为0。

# 6、以下编译能通过吗？
```go
package main

import (
	"fmt"
)

type People interface {
	Speak(string) string
}

type Stduent struct{}

func (stu *Stduent) Speak(think string) (talk string) {
	if think == "bitch" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}

func main() {
	var peo People = Stduent{}
	think := "bitch"
	fmt.Println(peo.Speak(think))
}
```
**解析:**
> 因为`func (stu *Student)Speak(think string)(talk string)`是`*Student{}`的方法，`Studnt{}`并没有实现该方法。

**解法一：**
> 定义为指针`var peo People = &Student{}`

**解法二：**
> 方法的receiver不用指针，`func (stu Student)Speak(think string)(talk string)`

# 7、下面代码有什么问题？
```go
type UserAges struct {
	ages map[string]int
	sync.Mutex
}

func (ua *UserAges) Add(name string, age int) {
	ua.Lock()							//ua.RLock
	ua.ages=make(map[string]int)
	defer ua.Unlock()					//defer ua.RUnlock
	ua.ages[name] = age
}

func (ua *UserAges) Get(name string) int {
	if age, ok := ua.ages[name]; ok {
		return age
	}
	return -1
}
```
**解析:**
> 1、最大的一个问题：
map是为nil的。并没有初始化，而nil的map不能赋值。是map是为nil的。并没有初始化，而nil的map不能赋值。

```go
func (ua *UserAges) Add(name string, age int) {
	ua.RLock()
	ua.ages=make(map[string]int)
	defer ua.RUnlock()
	ua.ages[name] = age
}
```

> 2、第二个问题：
对于读没有加锁，可能会引发panic。将`sync.Mutex`改为读写锁`sync.RWMutex`。