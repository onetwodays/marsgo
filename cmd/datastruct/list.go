package datastruct
//链表的一些算法
//反转链表
func Reverse(node *LNode){
	if node==nil || node.Next==nil {
		return
	}
	var pre  *LNode
	var next *LNode
	var cur *LNode = node.Next //初始化为第一个节点
	for cur!=nil{
		next = cur.Next //下一个节点
		cur.Next = pre
		pre = cur
		cur =next
	}
	node.Next=pre
}
//头插法逆序
func ReverseHead(node *LNode){
	if node==nil || node.Next==nil {
		return
	}

	cur:=node.Next.Next //从第二个节点开始
	node.Next.Next=nil//第一个节点作为未结点
	for cur!=nil{
		next:=cur.Next
		cur.Next=node.Next
		node.Next=cur
		cur=next
	}
}

func FindMiddleNode(head *LNode) *LNode  {
	if head==nil || head.Next==nil{
		return head
	}
	fast:= head.Next
	slow:=head.Next
	for fast!=nil && fast.Next!=nil{ //这里 fast.Next!=nil 必须加上
		fast=fast.Next.Next
		slow=slow.Next
	}
	return slow

}
//中间节点分2个链表
func SpitMiddleNode(head *LNode) *LNode  {
	if head==nil || head.Next==nil{
		return head
	}
	fast:= head.Next
	slow:=head.Next
	slowPre:=head.Next
	for fast!=nil && fast.Next!=nil{ //这里 fast.Next!=nil 必须加上
		slowPre=slow
		fast=fast.Next.Next
		slow=slow.Next
	}
	slowPre.Next=nil //分割链表为2部分
	return slow

}

//链表倒数第K个节点
func FindLastK(head *LNode,k int) *LNode  {
	if head==nil || head.Next==nil{
		return head
	}
	fast:=head.Next
	slow:=head.Next
	i:=1
	for ;i<=k && fast!=nil;i++{
		fast=fast.Next
	}
	if i<=k{
		return nil

	}
	for fast!=nil{
		fast=fast.Next
		slow=slow.Next
	}
	return slow

}
//相遇第一点为环的入口点
//快慢指针相遇时不一定是环的入口点
//设链表长度是L，入口环到相遇点x，起点到环入口是a
//L-a-x 是相遇点到入口环的距离
//链表头到入口环(n-1)*环长+相遇点到入口环的距离
//a=(n-1)r+(L-a-x)
func FindLoopNode(head *LNode) (first *LNode,isLoop bool) {
	first=nil
	isLoop=false
	if head==nil || head.Next==nil{
		return first,isLoop
	}
	slow:=head.Next
	fast:=head.Next
	for fast!=nil && fast.Next!=nil{
		slow=slow.Next
		fast=fast.Next.Next
		if slow==fast{
			isLoop=true
		}
	}
	//从链表头和相遇点分别设一个指针，每次各走一步，2个指针必定相遇，且相遇第一点为环入口点
	if isLoop{
		first:=head.Next

		for first!=slow{
			first=first.Next
			slow=slow.Next
		}
	}
	return first,isLoop
}

//
func ReverseValue(head *LNode)  {
	if head==nil || head.Next==nil || head.Next.Next==nil{
		return
	}

	pre:=head
	cur:=head.Next
	next:=head
	for cur!=nil && cur.Next!=nil{
		next = cur.Next.Next
		pre.Next=cur.Next
		cur.Next.Next=cur
		cur.Next=next
		pre=cur
		cur=next

	}
}