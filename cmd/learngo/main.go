package main

import (
    "fmt"
    "reflect"
    "unsafe"
)

type user struct {
    name string
    age int
}


type Slice []int

func (A Slice)Append(value int) {
    A = append(A, value)
}



func main1() {

    mSlice := make(Slice, 10, 20)
    mSlice.Append(5)
    fmt.Println(mSlice)


    i:= 10
    ip:=&i

    var fp *float64 = (*float64)(unsafe.Pointer(ip))

    *fp = *fp * 3

    fmt.Println(i)


    u:=new(user)


    pName:=(*string)(unsafe.Pointer(u))
    *pName="张三"

    pAge:=(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(u))+unsafe.Offsetof(u.age)))
    *pAge = 20

    fmt.Println(*u)

    fmt.Println(unsafe.Sizeof(true))
    fmt.Println(unsafe.Sizeof(int8(0)))
    fmt.Println(unsafe.Sizeof(int16(10)))
    fmt.Println(unsafe.Sizeof(int32(10000000)))
    fmt.Println(unsafe.Sizeof(int64(10000000000000)))
    fmt.Println(unsafe.Sizeof(int(10000000000000000)))


    var b bool
    var i8 int8
    var i16 int16
    var i64 int64

    var f32 float32

    var s string

    var m map[string]string

    var p *int32

    fmt.Println(unsafe.Alignof(b))
    fmt.Println(unsafe.Alignof(i8))
    fmt.Println(unsafe.Alignof(i16))
    fmt.Println(unsafe.Alignof(i64))
    fmt.Println(unsafe.Alignof(f32))
    fmt.Println(unsafe.Alignof(s))
    fmt.Println(unsafe.Alignof(m))
    fmt.Println(unsafe.Alignof(p))


    /*
    func Sizeof(x ArbitraryType) uintptr
     */
}

type User struct{
    Name string
    Age int
}

func main()  {

    u:= User{"张三",20}
    t:=reflect.TypeOf(u)
    fmt.Println(t)

    v:=reflect.ValueOf(u)
    fmt.Println(v)

    fmt.Printf("%T\n",u)
    fmt.Printf("%v\n",u)

    u1:=v.Interface().(User)
    fmt.Println(u1)


    t1:=v.Type()
    fmt.Println("t1",t1)

    fmt.Println(t.Kind())


    for i:=0;i<t.NumField();i++ {
        fmt.Println(t.Field(i).Name)
    }

    for i:=0;i<t.NumMethod() ;i++  {
        fmt.Println(t.Method(i).Name)
    }


    x:=2
    vq:=reflect.ValueOf(&x)
    vq.Elem().SetInt(100)
    fmt.Println(x)

}

/*

任何指针都可以转换为unsafe.Pointer
unsafe.Pointer可以转换为任何指针
uintptr可以转换为unsafe.Pointer
unsafe.Pointer可以转换为uintptr

取对齐值还可以使用反射包的函数，也就是说：unsafe.Alignof(x)等价于reflect.TypeOf(x).Align()
字段的偏移量，就是该字段在struct结构体内存布局中的起始位置(内存位置索引从0开始)。
在分析之前，我们先看下内存对齐的规则：

对于具体类型来说，对齐值=min(编译器默认对齐值，类型大小Sizeof长度)。也就是在默认设置的对齐值和类型的内存占用大小之间，取最小值为该类型的对齐值。我的电脑默认是8，所以最大值不会超过8.
struct在每个字段都内存对齐之后，其本身也要进行对齐，对齐值=min(默认对齐值，字段最大类型长度)。这条也很好理解，struct的所有字段中，最大的那个类型的长度以及默认对齐值之间，取最小的那个。
以上这两条规则要好好理解，理解明白了才可以分析下面的struct结构体。在这里再次提醒，对齐值也叫对齐系数、对齐倍数，对齐模数。这就是说，每个字段在内存中的偏移量是对齐值的倍数即可。
 */