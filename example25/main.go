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

func main() {
	//学生1 到 复兴复印店 打印信息
	s1 := &Student{
		Show: &FuxingStore{},
		Name: "学生1",
		Age:  20,
	}

	info1 := s1.ShowInfomation()
	fmt.Printf(info1)

	//学生2 到 东方打字复印店 打印信息
	s2 := &Student{
		Show: &DongfangStore{},
		Name: "学生2",
		Age:  21,
	}

	info2 := s2.ShowInfomation()
	fmt.Printf(info2)
}
