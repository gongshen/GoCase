//1：使用切片来实现栈的操作
package stack

import (
	"fmt"
)

type Element interface{}

//切片实现
type sliceInplement struct {
	element []Element
}

func NewSliceInplement() *sliceInplement {
	return &sliceInplement{}
}

func (si *sliceInplement) Push(e Element) {
	si.element = append(si.element, e)
}

func (si *sliceInplement) Pop() (e Element) {
	size := si.Size()
	if size == 0 {
		fmt.Println("The stack is empty!")
		return nil
	}
	lastElement := si.element[size-1]
	si.element[size-1] = nil
	si.element = si.element[:size-1]
	return lastElement
}

func (si *sliceInplement) Top() (e Element) {
	size := si.Size()
	if size == 0 {
		fmt.Println("The stack is empty!")
		return nil
	}
	lastElement := si.element[size-1]
	return lastElement
}

func (si *sliceInplement) Clear()bool{
	if si.IsEmpty(){
		fmt.Println("The stack is empty!")
		return false
	}
	for i:=0;i<si.Size();i++{
		si.element[i]=nil
	}
	si.element=make([]Element,0)
	return true
}

func (si *sliceInplement) Size() int {
	return len(si.element)
}

func (si *sliceInplement) IsEmpty() bool {
	if len(si.element) == 0 {
		return true
	}
	return false
}

