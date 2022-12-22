package main

import (
	"fmt"
	"goBase/go-learning/1-basic/7_interface/inherit/util"
	"io"
	"strings"
)

//-----------------------------struct 使用嵌套的 struct 和 interface---------------------------------------------------

type M struct {
	num int
}

func (m *M) Add(n int) {
	m.num += n
}

type S1 struct {
	*M        //嵌套 struct
	io.Reader //嵌套 interface
}

type S2 struct {
	M //嵌套 struct
}

func structEmbedBothExample() {
	//1. 使用嵌套的 struct 和 interface
	r := strings.NewReader("hello, go")
	s := S1{
		M:      &M{num: 17},
		Reader: r,
	}

	var sl = make([]byte, len("hello, go"))
	s.Read(sl)
	fmt.Println(string(sl))
	s.Add(5)
	fmt.Println(s.num)

	//2.结构体类型的方法集合，包含嵌入的 struct 和 interface 的方法集合。
	var s1 S1
	util.DumpMethodSet(s1)
	//main.S's method set:
	//- Add
	//- Read

	var s2 S2
	util.DumpMethodSet(s2)
	//main.S2's method set is empty!
}

//----------------------------------struct 中嵌套 struct----------------------------------------------
//T1 与 *T1、T2 与 *T2 的方法集合：
//   T1 的方法集合包含：T1M1；
//  *T1 的方法集合包含：T1M1、PT1M2；
//   T2 的方法集合包含：T2M1；
//  *T2 的方法集合包含：T2M1、PT2M2。
//
//它们作为嵌入字段嵌入到 T 中后，对 T 和 *T 的方法集合的影响也是不同的：
//  类型 T 的方法集合 = T1 的方法集合 + *T2 的方法集合
//  类型 *T 的方法集合 = *T1 的方法集合 + *T2 的方法集合

type T1 struct{}

func (T1) T1M1()   { println("T1's M1") }
func (*T1) PT1M2() { println("PT1's M2") }

type T2 struct{}

func (T2) T2M1()   { println("T2's M1") }
func (*T2) PT2M2() { println("PT2's M2") }

type T struct {
	T1
	*T2
}

func structEmbedStruct() {
	t := T{
		T1: T1{},
		T2: &T2{},
	}

	util.DumpMethodSet(t)
	//main.T's method set:
	//- PT2M2
	//- T1M1
	//- T2M1
	util.DumpMethodSet(&t)
	//*main.T's method set:
	//- PT1M2
	//- PT2M2
	//- T1M1
	//- T2M1
}

//------------------------------defined和alias类型的 method set--------------------------------------------------
//defined 类型的 method set

type T3 struct{}

func (T3) T3M1()  {}
func (*T3) T3M2() {}

type DT3 T3

func defineExample() {
	var t T3
	var pt *T3
	var dt3 DT3
	var pdt3 *DT3

	util.DumpMethodSet(t)
	//main.T3's method set:
	//- T3M1
	util.DumpMethodSet(pt)
	//*main.T3's method set:
	//- T3M1
	//- T3M2

	util.DumpMethodSet(dt3)
	//main.DT3's method set is empty!
	util.DumpMethodSet(pdt3)
	//*main.DT3's method set is empty!
}

//type alias 的 method set

type AT3 = T3

func aliasExample() {
	var t T3
	var pt *T3
	var at3 AT3
	var pat3 *AT3

	util.DumpMethodSet(t)
	//main.T3's method set:
	//- T3M1
	util.DumpMethodSet(pt)
	//*main.T3's method set:
	//- T3M1
	//- T3M2
	util.DumpMethodSet(at3)
	//main.T3's method set:
	//- T3M1
	util.DumpMethodSet(pat3)
	//*main.T3's method set:
	//- T3M1
	//- T3M2
}

//------------------------------------struct 中嵌套 interface--------------------------------------------

type I1 interface {
	I1M1()
	I1M2()
}

type T4 struct {
	I1
}

func interfaceExample() {
	var t4 T4
	var pt4 *T4

	util.DumpMethodSet(t4)
	//main.T4's method set:
	//- I1M1
	//- I1M2
	util.DumpMethodSet(pt4)
	//*main.T4's method set:
	//- I1M1
	//- I1M2
}

func main() {
	structEmbedBothExample()
	structEmbedStruct()

	defineExample()
	aliasExample()

	interfaceExample()
}
