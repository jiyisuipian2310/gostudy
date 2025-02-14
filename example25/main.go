package main

import "fmt"

//打印信息基类
type IBase interface {
	ShowInfomation(name string, age int) string
}

//复兴复印店
type FuxingStore struct{}

func (s *FuxingStore) ShowInfomation(name string, age int) (out string) {
	return fmt.Sprintf("复兴复印店打印信息: \n    姓名: %s\n    年龄: %d\n", name, age)
}

//东方打字复印店
type DongfangStore struct{}

func (s *DongfangStore) ShowInfomation(name string, age int) (out string) {
	return fmt.Sprintf("东方打字复印店打印信息: \n    name: %s\n    age: %d\n", name, age)
}

type Student struct {
	Show IBase
	Name string
	Age  int
}

func (s *Student) ShowInfomation() (out string) {
	return s.Show.ShowInfomation(s.Name, s.Age)
}

//设计模式中的依赖倒置原则
//将依赖关系从代码中分离出来，通过将依赖关系定义为接口类型, 可以在运行时动态地替换实现
func main() {
	students := []*Student{
		&Student{
			Show: &FuxingStore{},
			Name: "学生1",
			Age:  20,
		},
		&Student{
			Show: &DongfangStore{},
			Name: "学生2",
			Age:  21,
		},
		&Student{
			Show: &DongfangStore{},
			Name: "学生3",
			Age:  34,
		},
	}

	for _, s := range students {
		info := s.ShowInfomation()
		fmt.Printf(info)
	}
}
