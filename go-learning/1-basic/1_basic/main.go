package main

import "fmt"

func main() {
	//Go 的数据类型分四大类：
	//• 基础类型(basic type)：包括数字 (number)、 字符串(string) 和布尔型(boolean) 。
	//• 聚合类型(aggregate type)：数组 (array)和结构体(struct)一一是通过组合各种简单类型得到的更复杂的数据类型。
	//• 引用类型(reference type)：其中包含多种不同类型，如指针(pointer)，slice，map，函数(function)，以及通道(channel)。
	//    它们的共同点是全都间接指向程序变量或状态，于是操作所引用数据的效果就会遍及该数据的全部引用。
	//• 接口类型(interface type)：
	//当声明变量的时候，会做默认的赋 0 初始化。每种数据类型的默认赋 0 初始化的 0 值不同，例如 int 类型的 0 值为数值 0，float 的 0 值为 0.0，
	//  string 类型的 0 值为空 ""，bool 类型的 0 值为 false，数据结构的 0 值为 nil，struct 的 0 值为字段全部赋 0。

	//1.变量和常量
	//声明初始化一个变量
	var s string = "hello world"
	fmt.Println(s)
	//声明初始化多个变量
	var a1, a2, a3 int = 1, 2, 3
	fmt.Println(a1, a2, a3)
	//不用指明类型，通过初始化值来推导
	var b = true //bool型
	fmt.Println(b)
	//赋值初始化
	x := 100 //等价于 var x int = 100;
	fmt.Println(x)

	// 常量
	const pi float32 = 3.1415926

	//2.指针
	var i int = 1
	var pInt *int = &i
	fmt.Printf("i=%d\tpInt=%p\t*pInt=%d\n", i, pInt, *pInt) //输出：i=1     pInt=0xf8400371b0       *pInt=1
	*pInt = 2
	fmt.Printf("i=%d\tpInt=%p\t*pInt=%d\n", i, pInt, *pInt) //输出：i=2     pInt=0xf8400371b0       *pInt=2
	i = 3
	fmt.Printf("i=%d\tpInt=%p\t*pInt=%d\n", i, pInt, *pInt) //输出：i=3     pInt=0xf8400371b0       *pInt=3

	//3.new 和 make 内存分配
	//new 是一个分配内存的内建函数，但不同于其他语言中同名的 new 所作的工作，它只是将内存清零，而不是初始化内存。new(T) 为一个类型为 T 的新项目
	//  分配了值为零的存储空间并返回其地址，也就是一个类型为 *T 的值。用 Go 的术语来说，就是它返回了一个指向新分配的类型为 T 的零值的指针。
	//make(T, args) 函数的目的与 new(T) 不同。它仅用于创建 slice、map 和 chan（消息管道），并返回类型 T（不是 *T）的一个被初始化了的（不是零）实例。
	//  这种差别的出现是由于这三种类型实质上是对在使用前必须进行初始化的数据结构的引用。
	//  例如，slice 是一个具有三项内容的描述符，包括指向数据（在一个数组内部）的指针、长度以及容量，在这三项内容被初始化之前，slice 值为 nil。
	//  对于 slice、map 和 channel，make 初始化了其内部的数据结构并准备了将要使用的值。

	// 不必要地使问题复杂化：
	var p *[]int = new([]int) // 为切片结构分配内存；*p == nil；很少使用
	fmt.Println(p)            //输出：&[]
	*p = make([]int, 10, 10)
	fmt.Println(p)       //输出：&[0 0 0 0 0 0 0 0 0 0]
	fmt.Println((*p)[2]) //输出： 0
	// 习惯用法:
	v := make([]int, 10) // 切片 v 现在是对一个新的有 10 个整数的数组的引用
	fmt.Println(v)       //输出：[0 0 0 0 0 0 0 0 0 0]

	//4.switch
	//注意：switch 语句没有 break，还可以使用逗号 case 多个值
	switch i {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	case 4, 5, 6:
		fmt.Println("four, five, six")
	default:
		fmt.Println("invalid value!")
	}
}
