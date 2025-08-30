package main

import (
	"fmt"
)

// 定义一个抽象的基类
type IBase interface {
	Eat() // 描述吃的行为
}

// 定义一个 Animal 类，成员变量 Name，成员函数 Eat
type Animal struct {
	Name string
}

func (animal *Animal) Eat() {
	fmt.Printf("%v is eating\n", animal.Name)
}

// 动物的构造函数
func newAnimal(name string) *Animal {
	return &Animal{
		Name: name,
	}
}

// 定义一个 Person 类，成员变量 Name，成员函数 Eat
type Person struct {
	Name string
}

func (person *Person) Eat() {
	fmt.Printf("%v is eating\n", person.Name)
}

// 人的构造函数
func newPerson(name string) *Person {
	return &Person{
		Name: name,
	}
}

func DoSomething(param IBase) {
	param.Eat()
}

func main() {
	animal := Animal{Name: "dog"}
	animal.Eat()

	newAnimal("cat").Eat()

	person := Person{Name: "Tom"}
	person.Eat()

	newPerson("Sammy").Eat()

	DoSomething(&animal)
	DoSomething(&person)
}
