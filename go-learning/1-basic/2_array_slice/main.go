package main

import "fmt"

// 注意：数组是值类型，赋值和传参会复制整个数组，因此只改变副本的值，不会改变数组本身的值，想要在函数内部修改数组的值，可通过指针来传递数组参数。
func arrayDouble(arr *[5]int) {
	for i, val := range arr {
		arr[i] = val * 2
	}
}

//Go 中函数的参数是按值传递的，而 slice 是引用类型，将 slice 传递给函数，实际上将 {ptr, capacity, len} 进行了拷贝，
//  所以和源 slice 指向同一个底层数组。换句话说，如果函数内部对 slice 进行了修改，有可能会直接影响函数外部的底层数组，从而影响其它 slice。
//  通常需要将原有 slice 作为函数返回值。
func sliceFoo(s []int) []int {
	for i, val := range s {
		s[i] = val * 2
	}
	for i := 0; i < 3; i++ {
		s = append(s, i)
	}
	return s
}

func main() {
	//1. 数组 -- 值类型
	//1.1 指定数组长度定义
	var a1 = [5]int{1, 2, 3, 4, 5} // 也可以写成 直接赋值初始化的方式：a1 := [5]int{1, 2, 3, 4, 5}
	for i := 0; i < len(a1); i++ {
		fmt.Println(i, a1[i])
	}

	//1.2 让编译器推导数组长度
	var a2 = [...]int{1, 2, 3, 4, 5}
	for index, value := range a2 {
		fmt.Println(index, value)
	}

	//1.3 在 Go 中，当一个变量被声明之后，都会立即对其进行默认的赋 0 初始化。
	//对 int 类型的变量会默认初始化为 0，对 string 类型的变量会初始化为空 ""，对布尔类型的变量会初始化为 false，对指针(引用)类型的变量会初始化为 nil。
	var a3 [5]*int
	fmt.Println(a3)

	//1.4 传递数组给函数
	var a4 = [5]int{1, 2, 3, 4, 5}
	arrayDouble(&a4)
	fmt.Println(a4)

	//2.切片 --- 引用类型
	//  因为数组的长度是固定的，因此在 Go 语言中很少直接使用数组。
	//slice 是一个拥有相同类型元素的可变长度的序列，它是基于数组类型做的一层封装，功能更灵活，支持自动扩容和收缩。
	//切片是一个引用类型，底层引用一个数组对象，它的内部结构包含指针址、长度和容量，指针指向第一个 slice 元素对应的底层数组元素的地址（slice
	//  的第一个元素并不一定就是数组的第一个元素），长度对应 slice 中元素的数目，容量一般是从 slice 的开始位置到底层数据的结尾位置，长度不能超过容量。
	//|---------------|---------|----------|
	//| 0xc00007dd70  |   3     |     5    |
	//|      ptr      | Length  | Capacity |
	//|---------------|---------|----------|
	//可以通过使用内置的 len() 函数求长度，使用内置的 cap() 函数求切片的容量。

	//2.1 切片的基本定义与初始化、遍历
	var s0 = []int{1, 2, 3}
	//s0 := []int{1, 2, 3}          //直接赋值初始化的方式创建 slice

	//遍历
	for i, val := range s0 {
		fmt.Println(i, val)
	}
	fmt.Println(len(s0), cap(s0), s0) //3 3 [1 2 3]

	//2.2 通过make创建切片,并用空值初始化（推荐）
	s1 := make([]int, 3, 5)
	fmt.Println(len(s1), cap(s1), s1) //3 5 [0 0 0]

	//2.3 nil slice 和 空slice
	//nil slice 表示它的指针为 nil，也就是这个 slice 不会指向哪个底层数组。也因此，nil slice 的长度和容量都为 0。
	//|--------|---------|----------|
	//|  nil   |   0     |     0    |
	//|  ptr   | Length  | Capacity |
	//|--------|---------|----------|
	var s2 []int // 声明一个 nil slice
	println(s2)  //[0/0]0x0
	//空 slice
	//虽然 nil slice 和 Empty slice 的长度和容量都为 0，输出时的结果都是 []，且都不存储任何数据，但它们是不同的。
	//   nil slice 不会指向底层数组，而空 slice 会指向底层数组，只不过这个底层数组暂时是空数组。
	//|--------|---------|----------|
	//|  ADDR  |   0     |     0    |
	//|  ptr   | Length  | Capacity |
	//|--------|---------|----------|
	s3 := make([]int, 0, 0)
	println(s3) //[0/0]0xc00007dce8

	//2.4 为切片添加元素
	s4 := make([]int, 0, 3)
	s4 = append(s4, 1)
	s4 = append(s4, 2)
	s4 = append(s4, 3)
	fmt.Println(len(s4), cap(s4), s4) //3 3 [1 2 3]
	//当slice的length已经等于capacity的时候，再使用append()给slice追加元素，会自动扩展底层数组的长度。
	s44 := append(s4, 4)
	fmt.Println(len(s4), cap(s4), s4)    // 3 3 [1 2 3] --- 原有 s4 的 len 和 capacity 不变，实际上底层数组的长度和 capacity 变化了
	fmt.Println(len(s44), cap(s44), s44) // 4 6 [1 2 3 4]  --- s44 的 len 和 capacity 变化了，

	//2.5 对 slice 进行切片 --- 进行的是浅拷贝
	//SLICE[A:B] 截取时"左闭右开"，提取 [A, B) 区间元素，新的 capacity = 原 capacity - A
	// SLICE[A:]  // 从 A 切到最尾部
	// SLICE[:B]  // 从最开头切到 B(不包含 B)
	// SLICE[:]   // 从头切到尾，等价于复制整个 SLICE
	s5 := []int{1, 2, 3, 4, 5}
	s6 := s5[1:3]
	println(s5)     // [5/5]0xc00007dd70
	fmt.Println(s5) // [1 2 3 4 5]
	println(s6)     // [2/4]0xc00007dd78
	fmt.Println(s6) // [2 3]

	//SLICE[A:B:C] 提取 [A, B) 区间元素，C - A 代表的是capacity
	s7 := []int{1, 2, 3, 4, 5}
	s8 := s7[1:3:3]
	println(s7)     // [5/5]0xc00001c180
	fmt.Println(s7) // [1 2 3 4 5]
	println(s8)     // [2/2]0xc00001c188
	fmt.Println(s8) // [2 3]

	//2.6 切片浅拷贝
	//浅拷贝情况1
	s10 := s5
	s10[0] = 5
	fmt.Println(s5) //[5 2 3]

	//浅拷贝情况2
	s11 := []int{1, 2, 3}
	s12 := s11[1:]
	s12[0] = 5
	fmt.Println(s11) //[1 5 3]
	fmt.Println(s12) //[5 3]

	//2.7 切片深拷贝
	s13 := []int{1, 2, 3}
	s14 := make([]int, 2, 5)
	copy(s14, s13)   //注意：只会拷贝到目标切片的 len 长度，超过的丢弃
	fmt.Println(s13) //[1 2 3]
	fmt.Println(s14) //[1 2]

	//2.8 切片的删除
	s15 := []int{1, 2, 3, 4, 5}
	s15 = append(s15[:2], s15[3:]...)
	fmt.Println(s15) //[1 2 4 5]

	//2.9 传递slice给函数
	s16 := []int{1, 2, 3}
	fmt.Println(s16) //[1 2 3]
	println(s16)     //[3/3]0xc0000be0c0
	s16 = sliceFoo(s16)
	fmt.Println(s16) //[2 4 6 0 1 2]
	println(s16)     //[6/6]0xc0000c0150
}
