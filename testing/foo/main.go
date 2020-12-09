package main

import "fmt"

func init() {
	var f Foo
	f.Value = 32
	fmt.Println("Hello foo!", f)
}
