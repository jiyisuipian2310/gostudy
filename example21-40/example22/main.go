package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string `json:"name"`
	Age  int
}

func testEmptyInterface(v interface{}) {
	switch v.(type) {
	case nil:
		fmt.Printf("param is nil type\n")
	case int:
		fmt.Printf("param is int type\n")
	case int8:
		fmt.Printf("param is int8 type\n")
	case uint8:
		fmt.Printf("param is uint8 type\n")
	case int16:
		fmt.Printf("param is int16 type\n")
	case uint16:
		fmt.Printf("param is uint16 type\n")
	case int32:
		fmt.Printf("param is int32 type")
	case uint32:
		fmt.Printf("param is uint32 type\n")
	case int64:
		fmt.Printf("param is int64 type\n")
	case uint64:
		fmt.Printf("param is uint64 type\n")
	case string:
		fmt.Printf("param is string type\n")
	case []int:
		fmt.Printf("param is []int type\n")
	case []string:
		fmt.Printf("param is []string type\n")
	default:
		typ := reflect.TypeOf(v)
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Struct:
			fmt.Printf("param(%#v) is struct type\n", v)
			for i := 0; i < val.NumField(); i++ {
				fmt.Printf("%v --> %v --> %v\n", i, typ.Field(i).Name, val.Field(i))
			}
		case reflect.Slice:
			fmt.Printf("param(%#v) is Slice type\n", v)
		default:
			panic(fmt.Sprintf("marshal(%#v): cannot handle type %T", v, v))
		}
	}
}

func main() {
	testEmptyInterface(10)
	testEmptyInterface("zhangsan")
	testEmptyInterface([]int{10, 20, 30})
	testEmptyInterface([]string{"aa", "bb", "cc"})
	testEmptyInterface(Person{Name: "zhansgan", Age: 34})
}
