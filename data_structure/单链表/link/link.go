package link

import (
	"fmt"
)

type Node struct{
	Data int
	Next *Node
}

func (head *Node)Create() *Node{
	head = nil
	return head
}

func (p *Node)PrintLink(){
	for p!=nil{
		fmt.Printf("{Data:%d,Next:%p,Ptr:%p}\n",p.Data,p.Next,p)
		p=p.Next
	}
}

func (newNode *Node)Insert(head *Node)*Node{
	var p0,p1 *Node
	p0=newNode
	p1=head
			if head==nil{
				head=p0
				head.Next=nil
	}else{

		for p1.Next != nil {
			p1 = p1.Next
		}
			p1.Next = p0
			p0.Next = nil

	}
	return head
}

func (delNode *Node)Delete(head *Node)*Node  {
	var p1,p2 *Node
	if head == nil{
		fmt.Println("List is nil!")
		return head
	}
	p1 = head
	for p1.Data!=delNode.Data && p1.Next!=nil{
		p2=p1
		p1=p1.Next
	}
	if p1.Data == delNode.Data{
		if p1==head{
			head=p1.Next
		}else{
			p2.Next=p1.Next
		}
	}else{
		fmt.Println("The Node is not found!")
	}
	return head
}