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
}

func main()  {
	Create()
}