package datastruct

import (
	"fmt"
	"testing"
)

func TestReverse(t *testing.T) {
	head:=&LNode{}
	CreateNode(head,8)
	PrintNode("逆序前:",head)
	Reverse(head)
	PrintNode("逆序后:",head)
}

func TestReverseHead(t *testing.T) {
	fmt.Println("............")
	head:=&LNode{}
	CreateNode(head,8)
	PrintNode("逆序前:",head)
	ReverseHead(head)
	PrintNode("逆序后:",head)
}

func TestFindMiddleNode(t *testing.T) {
	fmt.Println("............")
	head:=&LNode{}
	CreateNode(head,9)
	PrintNode("逆序前:",head)
	middle:=FindMiddleNode(head)
	fmt.Println(middle.Data)
	PrintNode("逆序后:",head)
}

func TestSpitMiddleNode(t *testing.T) {
	fmt.Println("............")
	head:=&LNode{}
	CreateNode(head,9)
	PrintNode("逆序前:",head)
	middle:=SpitMiddleNode(head)

	fmt.Println(middle.Data)
	PrintNode("分割后1:",head)
	l:=NewLNode()
	l.Next=middle
	PrintNode("分割后2:",l)
	ReverseHead(l)
	PrintNode("逆序2:",l)
	first:=head.Next
	second:=l.Next

	for first!=nil {
		tmp:= first.Next
		first.Next=second
		first=tmp

		//第一个到达末尾了
		if first!=nil{
			tmp=second.Next
			second.Next=first
			second=tmp
		}




	}

	//在first=nil退出循环
	//second.Next=tmp2

	PrintNode("合并:",head)
}

func TestFindLastK(t *testing.T) {
	fmt.Println("............")
	head:=&LNode{}
	CreateNode(head,9)
	PrintNode("逆序前:",head)

	k:=FindLastK(head,2)
	fmt.Println(k.Data)

}

func TestFindLoopNode(t *testing.T) {

	head:=&LNode{}
	CreateNode(head,9)
	PrintNode("逆序前:",head)
	second:=head.Next.Next
	rear:=second.Next.Next.Next.Next.Next.Next
	fmt.Println("最后一个元素：",rear.Data)

	rear.Next=second
	fmt.Println("loop元素：",rear.Next.Data)

	node,isLoop:=FindLoopNode(head)
	if isLoop{
		fmt.Println("是环，入口环是：",node.Data)
	}else {
		fmt.Println("不是环，")
	}




}

func TestReverseValue(t *testing.T) {

	head:=&LNode{}
	CreateNode(head,8)
	PrintNode("交换前:",head)
	ReverseValue(head)
	PrintNode("交换后:",head)




}