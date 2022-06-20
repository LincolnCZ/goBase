package main

import (
	"fmt"
)

type Point struct {
	X float64
	Y float64
}

// ScaleBy
//只要receiver是值类型的，无论是使用值类型的实例还是指针类型的实例，都是拷贝整个底层数据结构的，方法内部访问的和修改的都是实例的副本。
//   所以，如果有修改操作，不会影响外部原始实例。
//func (p Point) ScaleBy(factor float64) {
//	p.X *= factor
//	p.Y *= factor
//}
//只要receiver是指针类型的，无论是使用值类型的实例还是指针类型的实例，都是拷贝指针，方法内部访问的和修改的都是原始的实例数据结构。所以，如果
//   有修改操作，会影响外部原始实例。
func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor
}

type Circle struct {
	Radius float64
	P      *Point // 具名嵌套struct
}

type Rectangle struct {
	*Point // 匿名嵌套struct
	Width  float64
	Height float64
}

//Go函数给参数传递值的时候是以复制的方式进行的，所以为了避免复制的开销，以及函数中的修改对原有值有效，大多数时候，传递给函数的数据结构都是它们的指针
func foo(device *Point) {
	device.X += 1
	device.Y += 1
}

func main() {
	//1.结构体的实例化、赋值
	//普通赋值方法
	var s1 Point
	s1.X = 1
	s1.Y = 1
	fmt.Printf("%#v\n", s1) //main.Point{X:1, Y:1}

	//实例化同时赋值1
	s2 := Point{
		X: 1,
		Y: 1,
	}
	fmt.Printf("%#v\n", s2) //main.Point{X:1, Y:1}

	//实例化同时赋值2
	s3 := Point{1, 1}
	fmt.Printf("%#v\n", s3) //main.Point{X:1, Y:1}

	//使用new 获取struct类型地址
	var s4 = new(Point)
	fmt.Printf("%#v\n", s4) //&main.Point{X:0, Y:0}
	//使用&x{} 初始化并获取地址
	s5 := &Point{
		X: 1,
		Y: 1,
	}
	fmt.Printf("%#v\n", s5) //&main.Point{X:1, Y:1}

	//2. nil struct
	//结构体的零值由结构体成员的零值组成。

	//3. 传递给函数
	s6 := &Point{
		X: 1,
		Y: 1,
	}
	fmt.Printf("%#v\n", s6) //&main.Point{X:1, Y:1}
	foo(s6)
	fmt.Printf("%#v\n", s6) //&main.Point{X:2, Y:2}

	//4.结构体嵌套
	//

	//具名嵌套结构体赋值方式、调用其嵌套struct的成员、函数
	s7 := Circle{
		Radius: 1,
		P: &Point{
			X: 1,
			Y: 1,
		},
	}
	fmt.Printf("%#v\n", s7) //main.Circle{Radius:1, P:(*main.Point)(0xc00012c170)}
	s7.P.ScaleBy(2)
	fmt.Printf("x:%f, y:%f \n", s7.P.X, s7.P.Y) // x:2.000000, y:2.000000

	//匿名嵌套struct、调用其嵌套struct的成员、函数
	s8 := Rectangle{
		Point: &Point{
			X: 1,
			Y: 1,
		},
		Width:  2,
		Height: 3,
	}
	fmt.Printf("%#v\n", s8) // main.Rectangle{Point:(*main.Point)(0xc00012c1a0)
	s8.ScaleBy(2)
	fmt.Printf("x:%f, y:%f \n", s8.X, s8.Y) // x:2.000000, y:2.000000
}
