package main

import (
	"encoding/binary"
	"fmt"
)

type MyStruct struct {
	Field1 int32
	Field2 string
	Field3 []int16
}

func main() {
	number := 1230
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b[0:], uint32(number))
	H := fmt.Sprintf("%x", b)
	fmt.Println(H)

	var s1 = MyStruct{123, "456", []int16{1, 2, 3}}
	var s2 MyStruct

	s2.Unmarshal(s1.Marshal())

	fmt.Println(s1)
	fmt.Println(s2)
}

func (s *MyStruct) binarySize() int {
	return 4 + // Field1
		2 + len(s.Field2) + // Len + Field2
		2 + 2*len(s.Field3) // Len + Field3
}

func (s *MyStruct) Marshal() []byte {
	b := make([]byte, s.binarySize())
	n := 0

	binary.BigEndian.PutUint32(b[n:], uint32(s.Field1))
	n += 4

	binary.BigEndian.PutUint16(b[n:], uint16(len(s.Field2)))
	n += 2

	copy(b[n:], s.Field2)
	n += len(s.Field2)

	binary.BigEndian.PutUint16(b[n:], uint16(len(s.Field3)))
	n += 2

	for i := 0; i < len(s.Field3); i++ {
		binary.BigEndian.PutUint16(b[n:], uint16(s.Field3[i]))
		n += 2
	}

	return b
}

func (s *MyStruct) Unmarshal(b []byte) {
	n := 0

	s.Field1 = int32(binary.BigEndian.Uint32(b[n:]))
	n += 4

	x := int(binary.BigEndian.Uint16(b[n:]))
	n += 2

	s.Field2 = string(b[n : n+x])
	n += x

	s.Field3 = make([]int16, binary.BigEndian.Uint16(b[n:]))
	n += 2

	for i := 0; i < len(s.Field3); i++ {
		s.Field3[i] = int16(binary.BigEndian.Uint16(b[n:]))
		n += 2
	}
}
