package datastruct
//二进制数中1的个数
//右移位
func CountOne1(n int) int  {
	count:=0
	for n>0{
		if n&1==1{ //判断最后一位是不是1
			count++
		}
		n>>=1
	}
	return count //O(n) n是二进制位数
}
func CountOne2(n int)  int {
	count:=0
	for n>0{
		n=n&(n-1)
		count++

	}
	return count //O(m) m是1的个数
}

//判断1024！末尾有多少0
func ZeroCount(n int ) int   {
	count:=0
	for n>0{
		n=n/5
		count+=n
	}
	return count

}


