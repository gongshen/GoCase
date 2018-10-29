package stack

type Stack interface {
	Pop() Element
	Push(e Element)
	Clear() bool
	Top() Element
	Size() int
	IsEmpty() bool
}