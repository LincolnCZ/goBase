package main

import (
	"fmt"
	"reflect"
)

//反射是指在程序运行期对程序本身进行访问和修改的能力。程序在编译时，变量被转换为内存地址，变量名不会被编译器写入到可执行部分。
//   在运行程序时，程序无法获取自身的信息。
//支持反射的语言可以在程序编译期将变量的反射信息，如字段名称、类型信息、结构体信息等整合到可执行文件中，并给程序提供接口访问反射信息，
//   这样就可以在程序运行期获取类型的反射信息，并且有能力修改它们。
//Go 语言由 reflect 包提供的反射功能，它定义了两个重要的类型, Type 和 Value， 一个 Type 表示一个 Go 类型。任意接口值在反射中都
//   可以理解为由 reflect.Type 和 reflect.Value 两部分组成，并且 reflect 包提供了 reflect.TypeOf 和 reflect.ValueOf 两个
//   函数来获取任意对象的 Value 和 Type。

type UserId int

func main() {
	//1.TypeOf 返回的是 reflect.Type 类型
	//1.1 查看 TypeOf 返回
	//   函数 reflect.TypeOf 接受任意的 interface{} 类型, 并以 reflect.Type 形式返回其动态类型
	var a int = 64
	t1 := reflect.TypeOf(a)            //注意 TypeOf 的参数类型是 interface{}
	fmt.Printf("type:%v\n", t1.Name()) //type:int

	//1.2 区分 Type 返回值的 Name 和 Kind
	//   在反射中关于类型还划分为两种：类型（Type）和种类（Kind）。因为在 Go 语言中我们可以使用 type 关键字构造很多自定义类型，而种类（Kind）
	//   就是指底层的类型，但在反射中，当需要区分指针、结构体等大品种的类型时，就会用到种类（Kind）
	var id UserId = 10
	t2 := reflect.TypeOf(id)
	fmt.Printf("type:%v kind:%v\n", t2.Name(), t2.Kind()) //type:UserId kind:int

	//2.ValueOf 返回的是 reflect.Value 类型
	//2.1 查看 ValueOf 返回
	//   reflect.ValueOf() 返回的是 reflect.Value 类型，其中包含了原始值的值信息。reflect.Value 与原始值之间可以互相转换。
	//reflect.Value类型提供的获取原始值的方法如下：
	//• Interface() interface {} 将值以 interface{} 类型返回，可以通过类型断言转换为指定类型
	//• Int() int64 将值以 int 类型返回，所有有符号整型均可以此方式返回
	//• Uint() uint64 将值以 uint 类型返回，所有无符号整型均可以此方式返回
	//• Float() float64 将值以双精度（float64）类型返回，所有浮点数（float32、float64）均可以此方式返回
	//• Bool() bool 将值以 bool 类型返回
	//• Bytes() []bytes 将值以字节数组 []bytes 类型返回
	//• String() string 将值以字符串类型返回
	var n1 int = 100
	v := reflect.ValueOf(n1)
	fmt.Println(v.Int())

	//2.2 修改 ValueOf 返回
	//   如果想要在函数中通过反射修改变量的值，需要注意函数参数传递的是值拷贝，必须传递变量地址才能修改变量值。而反射中使用专有的 Elem() 方法来获取指针对应的值。
	var n2 int = 100
	v2 := reflect.ValueOf(&n2)
	fmt.Println(v2) //得到的是一个地址：0xc0000b4048
	v2.Elem().SetInt(200)
	fmt.Println(n2) //输出200，修改成功！
}
