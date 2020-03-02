## 条件编译
  >我们发现，条件编译的关键在于-tags=jsoniter,也就是-tags这个标志，
   这就是Go语言为我们提供的条件编译的方式之一。
  
  好了，回过头来看我们刚开始时json/json.go、json/jsoniter.go这两个Go文件的顶部，都有一行注释：

  // +build !jsoniter
  
  // +build jsoniter
  
  这两行是Go语言条件编译的关键。+build可以理解为条件编译tags的声明关键字，后面跟着tags的条件。
  
  // +build !jsoniter表示，tags不是jsoniter的时候编译这个Go文件。 // +build jsoniter表示，tags是jsoniter的时候编译这个Go文件。
  
  也就是说，这两种条件是互斥的，只有当tags=jsoniter的时候，才会使用json-iterator，其他情况使用encoding/json。
  
  小结
  利用条件编译，我们实现了灵活选择json解析库的目的，除此之外，有时间我再细讲，而且tags只是其中的一部分，Go语言还可以根据Go文件后缀进行条件编译

** 编译构建的话，改为go build -tags=jsoniter .即可 **