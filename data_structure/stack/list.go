//2：使用链表来实现栈的操作
package stack

import "fmt"

var header *entry

//双向链的长度：栈的长度
var size int

type entry struct {
	prev    *entry
	next    *entry
	element Element
}

func newEntry(prev, next *entry, e Element) *entry {
	return &entry{prev, next, e}
}

//如果这里使用了 header:=newEntry() 的话会产生空指针引用的错误
//因为变量header不在该函数的作用域，所以会重新定义新的变量header
//这样引用一个未分配内存的空指针变量会panic
func NewHeader() *entry {
	header = newEntry(nil, nil, nil)
	header.next = header
	header.prev = header
	return header
}

//将新节点插入在表头的前面
func addBefore(e *entry, element Element) Element {
	newEntry := newEntry(e.prev, e, element)
	newEntry.prev.next = newEntry
	newEntry.next.prev = newEntry
	size++
	return newEntry
}

func (e *entry) Push(element Element) {
	addBefore(e, element)
}

func (e *entry) Pop() Element {
	if e.IsEmpty() {
		fmt.Println("The stack is empty!")
		return nil
	}
	prevEntry := header.prev
	header.prev = prevEntry.prev
	prevEntry.prev.next = header
	size--
	return prevEntry.element
}

func (e *entry) Clear() bool {
	if e.IsEmpty() {
		fmt.Println("The Stack is empty!")
		return false
	}
	list := header.next
	for list != header {
		next_list := list.next
		list.next = nil
		list.element = nil
		list.prev = nil
		list = next_list
	}
	header.next = header
	header.prev = header
	size = 0
	return true
}

func (e *entry) Top() Element {
	if e.IsEmpty() {
		fmt.Println("The Stack is empty!")
		return nil
	}
	return header.prev.element
}

func (e *entry) Size() int {
	return size
}

func (e *entry) IsEmpty() bool {
	if size == 0 {
		return true
	}
	return false
}
