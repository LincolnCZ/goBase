package main

import (
	"fmt"
	"reflect"
)

func DumpMethodSet(i interface{}) {
	v := reflect.TypeOf(i)
	elemTyp := v.Elem()

	n := elemTyp.NumMethod()
	if n == 0 {
		fmt.Printf("%s's method set is empty!\n", elemTyp)
		return
	}

	fmt.Printf("%s's method set:\n", elemTyp)
	for j := 0; j < n; j++ {
		fmt.Println("-", elemTyp.Method(j).Name)
	}
	fmt.Printf("\n")
}

type Interface interface {
	M1() float64
	M2() float64
}

type Test struct{}

func (t Test) M1() float64  { return 0.0 }
func (t *Test) M2() float64 { return 0.0 }

type Shape interface {
	M1() float64 //计算面积
	M2() float64 //计算周长
}

type R struct {
}

func (r *R) M1() float64 { //面积
	return 0.0
}

func (r *R) M2() float64 { //周长
	return 0.0
}

func main() {
	var t R
	var pt *R
	DumpMethodSet(&t)
	DumpMethodSet(&pt)
	//DumpMethodSet((*Interface)(nil))
}
