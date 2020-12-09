package main

import "fmt"

func init() {
	var b Bar
	b.Value = 32
	fmt.Println("Hello foo!", b)
}
