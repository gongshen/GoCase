
```go
type Person struct {
	Name string
	age int
}

func (p Person)String()string  {
	return fmt.Sprintf("%s:%d",p.Name,p.age)
}

type ByAge []Person

func (b ByAge)Len()int{return len(b)}
func (b ByAge)Swap(i,j int){b[i],b[j]=b[j],b[i]}
func (b ByAge)Less(i,j int)bool{return b[i].age<b[j].age}

func main()  {
	people:=[]Person{
		{"Bob",20},
		{"Tom",19},
		{"Miracle",21},
	}
	fmt.Printf("%T\n",ByAge{})
	sort.Sort(ByAge(people))
	fmt.Println(people)
}
```
在上面的例子中，Sort方法的参数是Interface,ByAge实现了Interface的3中方法，就等于实现了Interface。
```go
type Interface interface {
        Len() int
        Less(i, j int) bool
        Swap(i, j int)
}
```