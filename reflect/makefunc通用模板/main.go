package main

import (
	"fmt"
	"reflect"
)

func main() {
	//swap是传递给makefunc工作的
	swap := func(in []reflect.Value) []reflect.Value {
		return []reflect.Value{in[1], in[0]}
	}
	makeSwap := func(fptr interface{}) {
		//获取函数本身
		fn := reflect.ValueOf(fptr).Elem()
		v := reflect.MakeFunc(fn.Type(), swap)
		//分配给fn显示v
		fn.Set(v)
	}

	// Make and call a swap function for ints.
	var intSwap func(int, int) (int, int)
	makeSwap(&intSwap)
	//将参数0，1转换成值，调用swap
	fmt.Println(intSwap(0, 1))

	// Make and call a swap function for float64s.
	var floatSwap func(float64, float64) (float64, float64)
	makeSwap(&floatSwap)
	fmt.Println(floatSwap(2.72, 3.14))

}