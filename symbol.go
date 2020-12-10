package plugins

import "unsafe"

// Symbol represents a plugin symbol
type Symbol uintptr

// AsEmptyFunc will return the Symbol as an empty function
func (s Symbol) AsEmptyFunc() (fn func()) {
	return *(*func())(unsafe.Pointer(&s))
}

// AsInterfaceFunc will return the Symbol as an error function
func (s Symbol) AsInterfaceFunc() (fn func() interface{}) {
	return *(*func() interface{})(unsafe.Pointer(&s))
}

// AsErrorFunc will return the Symbol as an error function
func (s Symbol) AsErrorFunc() (fn func() error) {
	return *(*func() error)(unsafe.Pointer(&s))
}
