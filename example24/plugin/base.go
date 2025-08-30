package plugin

import "fmt"

type Plugin interface {
	Execute(args map[string]interface{}) (map[string]interface{}, error)
}

func init() {
	fmt.Printf("base.go init ...\n")
}
