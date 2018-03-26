# 将函数当作参数，进行二元操作
```go
type binaryOperation func(int,int)(int,error)

func operate(op1 int,op2 int,op binaryOperation)(result int,err error){
	if op==nil{
		err=errors.New("invaild functions!")
		return
	}
	return op(op1,op2)
}

//用户可以自己实现除法操作
func divide(op1 int,op2 int)(result int,err error){
	if op2==0{
		err=errors.New("divide by zero!")
		return
	}
	result=op1/op2
	return
}

func main()  {
	var result int
	var err error
	result,err=operate(0,2,divide)
	fmt.Printf("result：%v，err：%v\n",result,err)
}
```