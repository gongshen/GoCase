package main

import (
	"./link"
)

func main()  {
	var head *link.Node
	data1:=&link.Node{1,nil}
	data2:=&link.Node{2,nil}
	data3:=&link.Node{3,nil}
	data4:=&link.Node{4,nil}
	head=data1.Insert(head)
	head=data2.Insert(head)
	head=data3.Insert(head)
	head=data4.Insert(head)

	head.PrintLink()
	head=data3.Delete(head)

	head.PrintLink()
}
