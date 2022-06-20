package main

import (
	"fmt"
	"goBase/go-learning/1-basic/7_interface/inherit/util"
)

//-----------------------------struct 中嵌套 *struct 达到“继承”---------------------------------------------------

// IAnimal 模拟动物行为的接口
type IAnimal interface {
	Eat() // 描述吃的行为
}

// Animal 所有动物的父类
type Animal struct {
	Name string
}

// Eat 动物去实现IAnimal中描述的吃的接口
func (a *Animal) Eat() {
	fmt.Printf("%v is eating\n", a.Name)
}

// NewAnimal 动物的构造函数
func NewAnimal(name string) *Animal {
	return &Animal{
		Name: name,
	}
}

// Cat 组合了animal
type Cat struct {
	*Animal
	//Animal //使用这种方式也可以达到继承的效果？
}

// NewCat 实现猫的构造函数，初始化animal结构体
func NewCat(name string) *Cat {
	return &Cat{
		Animal: NewAnimal(name),
	}
}

func inheritExample() {
	cat := NewCat("cat")
	cat.Eat() // cat is eating
	util.DumpMethodSet(cat)
	//*main.Cat's method set:
	//- Eat
}

//-----------------------------struct 中嵌套 struct 达到“继承”---------------------------------------------------

type I2 interface {
	I2M1()
	I2M2()
}

type I3 interface {
	I3M1()
	I3M2()
}

type S2 struct {
}

func (s *S2) I2M1() {
}

func (s *S2) I2M2() {
}

type S3 struct {
}

func (s *S3) I3M1() {
}

func (s *S3) I3M2() {
}

// S 类型的实例，实现了 I3
// *S 类型的实例，实现了 I2、I3
type S struct {
	S2
	*S3
}

func inherit2Example() {
	var s S
	var ps *S

	util.DumpMethodSet(s)
	//main.S's method set:
	//- I3M1
	//- I3M2
	util.DumpMethodSet(ps)
	//*main.S's method set:
	//- I2M1
	//- I2M2
	//- I3M1
	//- I3M2

	var i2 I2
	i2 = ps
	util.DumpMethodSet(i2)
}

func main() {
	inheritExample()
	inherit2Example()
}
