package main

import (
	"fmt"
	"reflect"
)

func main() {
	u := User{1, "gongshen", 21}
	v := reflect.ValueOf(u)
	mv := v.MethodByName("Hello") //调用方法
	args := []reflect.Value{reflect.ValueOf("lutianqi")}
	mv.Call(args) //调用Call，传入slice
}

type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) Hello(name string) {
	fmt.Println("Hello", name, "my name is", u.Name)
}
