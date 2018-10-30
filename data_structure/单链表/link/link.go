package list

import (
	"fmt"
)

type Element interface {}

type node struct{
	Data Element
	Next *node
}

type SingalList interface {
	//Front(head *node,new *node)	//在链头插入
	Remove(head *node)SingalList	//删除
	//Back()：在链尾插入
	//Insert(head *node,index int,data Element)SingalList	//在指定位置插入
	Insert(head *node) SingalList
}

func (head *node)Create() SingalList{
	head = nil
	return head
}

func (p *node)PrintLink(){
	for p!=nil{
		fmt.Printf("{Data:%d,Next:%p,Ptr:%p}\n",p.Data,p.Next,p)
		p=p.Next
	}
}

func (newNode *node)Insert(head *node)SingalList{
	var p0,p1 *node
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

func (delNode *node)Remove(head *node)SingalList  {
	var p1,p2 *node
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