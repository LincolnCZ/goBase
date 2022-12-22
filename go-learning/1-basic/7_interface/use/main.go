package main

import (
	"fmt"
	"goBase/go-learning/1-basic/7_interface/inherit/util"
	"math"
)

//实现接口：一个接口类型定义了一套方法，如果一个具体类型要实现该接口，那么必须实现接口类型定义中的所有方法。如果一个类型实现了一个接口要求的所有
//  方法，那么这个类型实现了这个接口。为了简化表述，Go 程序员通常说一个具体类型 "是一个" (is-a) 特定的接口类型，这其实代表着该具体类型实现了该接口。
//接口的赋值规则：仅当一个表达式实现了一个接口时，这个表达式才可以赋给该接口。
//实例的 method set：实例的 method set 决定了它所实现的接口，以及通过 receiver 可以调用的方法。https://golang.org/ref/spec#Method_sets
// |实例的类型   | 包含的receiver（指的是struct方法中的receiver）方法|
// |------------|----------------------------------------------|
// |值类型:T     | (T Type)的方法                                |
// |指针类型:*T  | (T Type)或(T *Type)的方法                      |
//方法集合决定接口实现的含义就是：如果某类型 T 的方法集合与某接口类型的方法集合相同，或者类型 T 的方法集合是接口类型 I 方法集合的超集，那么我
//  们就说这个类型 T 实现了接口 I。或者说，方法集合这个概念在 Go 语言中的主要用途，就是用来判断某个类型是否实现了某个接口。

//如何选择 receiver 类型?
//（1）T 类型是否需要实现某个接口，也就是是否存在将 T 类型的变量赋值给某接口类型变量的情况。
//    如果 T 类型需要实现某个接口，那我们就要使用 T 作为 receiver 参数的类型，来满足接口类型方法集合中的所有方法。
//    如果 T 不需要实现某个接口，但 *T 需要实现该接口，那么根据方法集合概念，*T 的方法集合是包含 T 的方法集合的，这样我们在确定 Go 方法的
//      receiver 的类型时，参考原则二和原则三就可以了。
//（2）如果 Go 方法要把对 receiver 参数所代表的类型实例的修改反映到原类型实例上，那么我们应该选择 *T 作为 receiver 参数的类型。
//    否则通常我们会为 receiver 参数选择 T 类型，这样可以减少外部修改类型实例内部状态的“渠道”。
//（3）除非 receiver 参数类型的 size 较大，考虑到传值的较大性能开销，选择 *T 作为 receiver 类型可能更适合。
//    考虑到 Go 方法调用时，receiver 参数是以值拷贝的形式传入方法中的。那么，如果 receiver 参数类型的 size 较大，以值拷贝形式传入就会
//      导致较大的性能开销，这时我们选择 *T 作为 receiver 类型可能更好些。

type Shape interface {
	Area() float64      //计算面积
	Perimeter() float64 //计算周长
}

type Rect struct {
	width, height float64
}

func (r *Rect) Area() float64 { //面积
	return r.width * r.height
}

func (r *Rect) Perimeter() float64 { //周长
	return 2 * (r.width + r.height)
}

type Circle struct {
	radius float64
}

func (c *Circle) Area() float64 { //面积
	return math.Pi * c.radius * c.radius
}

func (c *Circle) Perimeter() float64 { //周长
	return 2 * math.Pi * c.radius
}

//将接口类型作为参数很常见。这时，那些实现接口的实例都能作为接口类型参数传递给函数。
func printArea(s Shape) {
	fmt.Printf("Area:%f, Perimeter:%f\n", s.Area(), s.Perimeter())
}

func main() {
	//1. 接口的使用
	r := Rect{width: 2.9, height: 4.8}
	c := Circle{radius: 4.3}

	//指针形式的 &r 才能赋值给interface；
	//值形式的 r 无法赋值给interface，因为 r 是值类型，对应的method set 不包含shape的接口实现。
	s := []Shape{&r, &c}
	r.Area()

	for _, sh := range s {
		fmt.Println(sh)
		fmt.Println(sh.Area())
		fmt.Println(sh.Perimeter())
	}

	//2.传递给函数
	printArea(&r) //需要使用指针形式 &r，r无法使用。理由同上

	//3. ...interface{} 作为函数参数
	//func Println(a ...interface{}) (n int, err error)
	//每一个参数都会放进一个名为a的Slice中，Slice中的元素是接口类型，而且是空接口，这使得无需实现任何方法，任何东西都可以丢
	//    到fmt.Println()中来，至于每个东西怎么输出，那就要看具体情况：由类型的实现的 String() 方法决定。

	//4.类型断言 x.(T) 检查 x 的动态类型是否是 T，其中 x 必须是接口值。

	//如果 T 是接口类型，类型断言检查 x 的动态类型是否满足 T。如果检查成功，x 的动态值不会被提取，返回值是一个类型为 T 的接口值。换句话说，
	//    到接口类型的类型断言，改变了表达式的类型，改变了（通常是扩大了）可以访问的方法，且保护了接口值内部的动态类型和值。
	var x interface{}
	x = &r
	v, ok := x.(Shape) //判断 x 是否实现了 shape interface
	if ok {
		fmt.Println("implement Shape interface", v) // implement Shape interface &{2.9 4.8}
	} else {
		fmt.Println("not implement Shape interface", v)
	}

	//如果 T 是具体类型，类型断言检查 x 的动态类型是否等于具体类型 T。如果检查成功，类型断言返回的结果是 x 的动态值，其类型是 T。换句话说，
	//    对接口值 x 断言其动态类型是具体类型 T，若成功则提取出 x 的具体值。如果检查失败则 panic。
	x = c
	v1, ok1 := x.(Circle)
	if ok1 {
		fmt.Println("the real type of x is Circle", v1) // the real type of x is Circle {4.3}
	} else {
		fmt.Println("the real type of x is not Circle", v1)
	}

	//也可以采用以下方式进行判断
	switch x.(type) {
	case *Circle:
		fmt.Println("is *Circle")
	case Circle:
		fmt.Println("is Circle")
	case *Rect:
		fmt.Println("is *Rect")
	case Rect:
		fmt.Println("is Rect")
	default:
		fmt.Println("unknown")
	}

	//5. 探索method set
	var t Circle
	util.DumpMethodSet(t)
	// main.Circle's method set is empty!

	var pt *Circle
	util.DumpMethodSet(pt)
	//*main.Circle's method set:
	//- Area
	//- Perimeter

	util.DumpMethodSet((*Shape)(nil))
	//*main.Shape's method set is empty!
}
