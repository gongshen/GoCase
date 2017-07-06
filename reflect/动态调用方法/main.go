package main

import (
	"fmt"
	"reflect"
)

func main() {
	u := User{1, "GS", 22}
	v := reflect.ValueOf(u)
	mv := v.MethodByName("Hello") //调用方法
	/*
	//第一种方法
	out:=mv.Call([]reflect.Value{
		reflect.ValueOf("%s=%d"),
		reflect.ValueOf("x"),
		reflect.ValueOf(20),
	})
	*/
	//第二种方法
	out:=mv.CallSlice([]reflect.Value{
		reflect.ValueOf("%s=%d"),
		reflect.ValueOf([]interface{}{"x",20}),
	})
	fmt.Println(out)
}

type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) Hello(str string,a ...interface{})string{
	return fmt.Sprintf(str,a...)
}