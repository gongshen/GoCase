package btree

import "fmt"

type Node struct {
	Left	*Node
	Data	interface{}
	Right	*Node
}
//初始化二叉树
type Initer interface {
	SetData(data interface{})
}
//对二叉树的操作
type Operate interface {
	PrintBT()				//打印
	Depth() 	int			//深度计算
	LeafCount()	int			//叶子统计
}
//遍历
type Order interface {
	PreOrder()
	InOrder()
	PostOrder()
}

func (n *Node)SetData(data interface{}){
	n.Data=data
}

func (n *Node)PrintBT()  {
	PrintBT(n)
}
func (n *Node)Depth()int  {
	return Depth(n)
}
func (n *Node)LeafCount()int {
	return LeafCount(n)
}
func (n *Node)PreOrder()  {
	PreOrder(n)
}
func (n *Node)InOrder()  {
	InOrder(n)
}
func (n *Node)PostOrder()  {
	PostOrder(n)
}

func NewNode(left,right *Node)*Node  {
	return &Node{left,nil,right}
}

func PrintBT(n *Node)  {
	if n!=nil{
		fmt.Printf("%v",n.Data)
		if n.Left!=nil || n.Right!=nil{
			fmt.Printf("(")
			PrintBT(n.Left)
			if n.Right!=nil{
				fmt.Printf(",")
			}
			PrintBT(n.Right)
			fmt.Printf(") ")
		}
	}
}

func Depth(n *Node) int  {
	var depLeft,depRight int
	if n == nil{
		return 0
	}else{
		depLeft=Depth(n.Left)
		depRight=Depth(n.Right)
		if depLeft>depRight{
			return depLeft+1
		}else{
			return depRight+1
		}
	}
}

func LeafCount(n *Node)int  {
	if n == nil{
		return 0
	}else if n.Left == nil && n.Right == nil{
			return 1
		}else{
		return LeafCount(n.Left)+LeafCount(n.Right)
	}
}

func PreOrder(n *Node)  {
	if n!=nil{
		fmt.Printf("Data:%v",n.Data)
		PreOrder(n.Left)
		PreOrder(n.Right)
	}
}

func InOrder(n *Node)  {
	if n!=nil{
		InOrder(n.Left)
		fmt.Printf("Data:%v",n.Data)
		InOrder(n.Right)
	}
}

func PostOrder(n *Node)  {
	if n!=nil{
		PostOrder(n.Left)
		PostOrder(n.Right)
		fmt.Printf("Data:%v",n.Data)
	}
}