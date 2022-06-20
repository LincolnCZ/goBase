package main

import "goBase/go-learning/1-basic/7_interface/inherit/util"

type Interface1 interface {
	Func1()
	Func2()
}

type T1 struct{}

func (t T1) Func1()  {}
func (t *T1) Func2() {}

func main() {
	//1. int、*int 为代表的 Go 原生类型由于没有定义方法，所以它们的方法集合都是空的。
	var n int
	util.DumpMethodSet(n)
	//int's method set is empty!
	util.DumpMethodSet(&n)
	//*int's method set is empty!

	//2.Go 语言规定，*T 类型的方法集合包含所有以 *T 为 receiver 参数类型的方法，以及所有以 T 为 receiver 参数类型的方法。
	var t T1
	var pt *T1
	util.DumpMethodSet(t)
	//main.T1's method set:
	//- Func1
	util.DumpMethodSet(pt)
	//*main.T1's method set:
	//- Func1
	//- Func2

	//3.
	var i Interface1
	i = pt
	util.DumpMethodSet(i)
	//*main.T1's method set:
	//- Func1
	//- Func2
	//i = t // 无法将 't' (类型 T1) 用作类型 Interface1 类型未实现 'Interface1'，因为 'Func2' 方法有指针接收器
}
