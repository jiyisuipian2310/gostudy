package main

import (
	"fmt"
)

/*
 * 定义一个Animal接口，包含Say()方法，该方法返回string类型。
 * 定义Dog和Cat两个结构体，实现Animal接口，并包含Say()方法。
 * 定义一个Say()函数，接收Animal接口作为参数，并根据Animal接口的类型，调用其Say()方法。
 */

type Animal interface {
	Say() string
	Run()
}

func Say(a Animal) string {
	switch t := a.(type) {
	case *Dog:
		fmt.Printf("This is a Dog class\n")
		return a.Say()
	case *Cat:
		fmt.Printf("This is a Cat class\n")
		return a.Say()
	default:
		return fmt.Sprintf("unknown type [%s]", t)
	}
}

/****************Dog class************************/
type Dog struct {
	Name string
	Age  int
}

func (d *Dog) Say() string {
	return fmt.Sprintf("I am a Dog, my name is %s, and I am %d years old", d.Name, d.Age)
}

func (d *Dog) Run() {
	fmt.Println(d.Name, " is running")
}

/****************Cat class************************/
type Cat struct {
	Name string
	Age  int
}

func (c *Cat) Say() string {
	return fmt.Sprintf("I am a Cat, my name is %s, and I am %d years old", c.Name, c.Age)
}

func (c *Cat) Run() {
	fmt.Println(c.Name, " is running")
}

func main() {
	dog := &Dog{"DogName", 20}
	content := Say(dog)
	fmt.Println(content)

	cat := &Cat{"CatName", 15}
	content = Say(cat)
	fmt.Println(content)
}
