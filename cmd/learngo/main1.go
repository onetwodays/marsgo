package main

import "fmt"

type Slice1 []int
func (A Slice1)Append(value int) {

	A1 := append(A, value)
	fmt.Println("函数内部：",A)
	fmt.Printf("%p\n%p\n",A,A1)
}


func main() {
	mSlice := make(Slice1, 10, 20)
	mSlice.Append(5)
	fmt.Println(mSlice)
}