package main

import (
	"reflect"
	"strings"
	"fmt"
)

func add(args []reflect.Value)(result []reflect.Value){
	if len(args)==0{
		return nil
	}
	var ret reflect.Value
	switch args[0].Kind() {
	case reflect.Int:
		n:=0
		for _,i:=range args{
			n+=int(i.Int())
		}
		ret=reflect.ValueOf(n)
	case reflect.String:
		var ss=make([]string,0,len(args))
		for _,s:=range args{
			ss=append(ss,s.String())
		}
		ret=reflect.ValueOf(strings.Join(ss,""))
	}
	result=append(result,ret)
	return
}

func makeAdd(ftpr interface{})  {
	fn:=reflect.ValueOf(ftpr).Elem()
	v:=reflect.MakeFunc(fn.Type(),add)
	fn.Set(v)
}

func main()  {
	var intAdd func(i,j int)int
	var strAdd func(a,b string)string
	makeAdd(&intAdd)
	makeAdd(&strAdd)
	fmt.Println(intAdd(100,200))
	fmt.Println(strAdd("Hello, ","world!"))
}
