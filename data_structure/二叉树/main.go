package main

import (
	"./btree"
)

func Create()  {
	root:=btree.NewNode(nil,nil)
	root.SetData("root node")
	leftNode :=btree.NewNode(nil,nil)
	leftNode.SetData("left node")
	leftChild:=btree.NewNode(nil,nil)
	leftChild.SetData(111)
	rightChild:=btree.NewNode(nil,nil)
	rightChild.SetData(3.1415926)
	leftNode.Left=leftChild
	leftNode.Right=rightChild

	rightNode:=btree.NewNode(nil,nil)
	rightNode.SetData("right node")

	root.Left=leftNode
	root.Right=rightNode
	root.PrintBT()
	fmt.Println()
}

func Operation()  {
	var op btree.Operate
	op=root
	fmt.Println("The Depth is:",op.Depth())
	fmt.Println("The nums of leaf:",op.LeafCount())
}

func Order()  {
	var or btree.Order
	or=root
	or.PreOrder()
	fmt.Println("先序：")
	or.InOrder()
	fmt.Println("中序：")
	or.PostOrder()
	fmt.Println("后序：")


func main()  {
	Create()
	Operation()
	Order()
}