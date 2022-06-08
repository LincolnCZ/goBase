package main

import "fmt"

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
}

// NewCat 实现猫的构造函数，初始化animal结构体
func NewCat(name string) *Cat {
	return &Cat{
		Animal: NewAnimal(name),
	}
}

func InheritExample() {
	cat := NewCat("cat")
	cat.Eat() // cat is eating
}
