package main

import "fmt"

//函数定义：
// func 函数名(参数)(返回值){
//    函数体
//}
//注意事项：
// 1. 函数的参数、返回值以及它们的类型，结合起来成为函数的签名(signature)。
// 2. 函数调用的时候，如果有参数传递给函数，则先拷贝参数的副本，再将副本传递给函数。
//  • 由于引用类型(slice、map、interface、channel)自身就是指针，所以这些类型的值拷贝给函数参数，函数内部的参数仍然指向它们的底层数据结构。
// 3. Go 中不允许函数重载(overload)，也就是说不允许函数同名。

//Greeting 变长参数
//  Go 语言中可通过在参数名后加 ... 来标识变长参数，变长参数是指函数的参数数量不固定，变长参数通常作为函数的最后一个参数，本质上，函数的变长参数是通过切片来实现的。
func Greeting(who ...string) {
	fmt.Printf("%#v\n", who) //[]string{"Joe", "Anna", "Eileen"}
	//在使用...的时候(如传递、赋值)，可以将变长参数看成是一个slice
	for _, val := range who {
		fmt.Printf("%s\n", val)
	}
}

//多返回值函数
//  返回值由返回值变量和其变量类型组成，也可以只写返回值的类型，Go 语言中函数支持多返回值，多个返回值必须用 () 包裹，并用英文逗号 , 分隔；
//  函数定义时可以给返回值命名，并在函数体中直接使用这些变量，最后通过 return 关键字返回。
func divide(a int, b int) (int, int) {
	var n1, n2 int
	n1 = a / b
	n2 = a % b
	return n1, n2
}

//函数 incr 返回了一个函数，返回的这个函数就是一个闭包。这个函数中本身是没有定义变量 x 的，而是引用了它所在的环境（函数 incr）中的变量 x。
func incr() func() int {
	var x int
	return func() int {
		x++
		return x
	}
}

//错误使用：捕获了迭代变量
//  变量 d 在 for 循环引进的一个块作用域内进行声明。在循环里创建的所有函数变量共享相同的变量--一个可访问的存储位置，而不是固定的值。
//  d 变量的值在不断地迭代中更新，因此当调用清理函数时，d 变量已经被每一次的 for 循环更新多次。因此，d 变量的实际取值是最后一次迭代时的值，所以输出结果均为 5。
func captureIteration() {
	s := []int{1, 2, 3, 4, 5}
	var printFuncs []func()
	for _, d := range s {
		printFuncs = append(printFuncs, func() { fmt.Println(d) })
	}
	for _, item := range printFuncs {
		item() // 输出的值均为5
	}
}

//第一种方法是在循环体内部再定义一个局部变量，这样每次迭代 printFuncs 语句的闭包函数捕获的都是不同的变量，这些变量的值对应迭代时的值。
func captureIterationFix() {
	s := []int{1, 2, 3, 4, 5}
	var printFuncs []func()
	for _, d := range s {
		d := d // 声明一个内部d，并已外部d初始化
		printFuncs = append(printFuncs, func() { fmt.Println(d) })
	}
	for _, item := range printFuncs {
		item()
	}
}

//第二种方式是将迭代变量通过匿名函数的参数传入，printFuncs 语句会马上对调用参数求值。
func captureIterationFix2() {
	s := []int{1, 2, 3, 4, 5}
	for _, d := range s {
		func(i int) {
			fmt.Println(i)
		}(d)
	}
}

func main() {
	//1.变长参数
	Greeting("Joe", "Anna", "Eileen")

	//2.多返回值函数
	n1, n2 := divide(21, 10)
	fmt.Println(n1, n2)

	//3.匿名函数
	//  匿名函数就是没有函数名的函数，当我们不希望给函数起名字的时候，可以使用匿名函数，匿名函数不能够独立存在，它可以被赋值于某个变量或直接对匿名函数进行调用。
	//3.1将匿名函数保存到变量
	add := func(x, y int) {
		fmt.Printf("The sum of %d and %d is: %d\n", x, y, x+y)
	}
	add(10, 20) // 通过变量调用匿名函数

	//3.2直接对匿名函数进行调用，最后的一对括号表示对该匿名函数的直接调用执行
	func(x, y int) {
		fmt.Printf("The sum of %d and %d is: %d\n", x, y, x+y)
	}(10, 20)

	//4. 闭包
	//  闭包指的是一个函数和与其相关的引用环境组合而成的实体(即：闭包 = 函数 + 引用环境)。

	//4.1 这里的 i 就成为了一个闭包，闭包对外层词法域变量是引用的，即 i 保存着对 x 的引用。
	i := incr()
	fmt.Println(i()) // 1
	fmt.Println(i()) // 2
	fmt.Println(i()) // 3

	//4.2 这里调用了三次 incr()，返回了三个闭包，这三个闭包引用着三个不同的 x，它们的状态是各自独立的
	fmt.Println(incr()()) // 1
	fmt.Println(incr()()) // 1
	fmt.Println(incr()()) // 1

	//4.3警告：捕获迭代变量--易错点
	captureIteration()
	captureIterationFix()
	captureIterationFix2()
}
