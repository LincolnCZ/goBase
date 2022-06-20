package main

type E interface {
	Func1()
	Func2()
}

// I 像这种在一个接口类型（I）定义中，嵌入另外一个接口类型（E）的方式，就是我们说的接口类型的类型嵌入。
//等价于：
//type I interface {
//	Func1()
//	Func2()
//	Func3()
//}
//这种接口类型嵌入的语义就是新接口类型（如接口类型 I）将嵌入的接口类型（如接口类型 E）的方法集合，并入到自己的方法集合中。
type I interface {
	E
	Func3()
}
