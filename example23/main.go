package main

import "fmt"

type tree struct {
	value       int
	left, right *tree
}

// Sort sorts values in place.
func Sort(values []int, reverse bool) {
	var root *tree
	for _, v := range values {
		root = add(root, v, reverse)
	}

	appendValues(values[:0], root)
}

func add(t *tree, value int, reverse bool) *tree {
	if t == nil {
		t = new(tree)
		t.value = value
		fmt.Printf("value: %d\n", value)
		return t
	}

	if reverse {
		if value > t.value {
			t.left = add(t.left, value, reverse)
		} else {
			t.right = add(t.right, value, reverse)
		}
	} else {
		if value < t.value {
			t.left = add(t.left, value, reverse)
		} else {
			t.right = add(t.right, value, reverse)
		}
	}

	return t
}

// appendValues appends the elements of t to values in order
// and returns the resulting slice.
func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func main() {
	var ages = []int{100, 10, 90, 20, 80, 30, 70, 40}
	fmt.Printf("ages: %v\n", ages)
	Sort(ages, false) //false 正向排序  true 方向排序
	fmt.Printf("ages: %v\n", ages)
}

