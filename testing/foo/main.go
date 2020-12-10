package main

import (
	"fmt"
	"math/big"
)

func main() {
	var f Foo
	f.Value = 32
	fmt.Println("Hello foobie!", f)
}

// BigInt will return a big.Int
func BigInt() interface{} {
	var i big.Int
	i.SetUint64(1337)
	return i
}
