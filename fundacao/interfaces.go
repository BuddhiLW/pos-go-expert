package fundacao

import "fmt"

// Interfaces
type I interface {
	F1()
}

func (s StruFoo) F1() {
	fmt.Printf("StruFoo prints itself: %v \n", s)
}

func (s StruFooBar) F1() {
	fmt.Printf("StruFooBar prints itself: %v \n", s)
}

func (s StruBarfoo) F1() {
	fmt.Printf("StruBarfoo print itself: %v \n", s)
}

func Interfaces() {
	struFoo := StruFoo{
		S: "Foo",
		I: 256,
		B: true,
	}

	struFooBar := StruFooBar{
		SS:      []string{"Foo", "Bar"},
		II:      []int{256, 512},
		BB:      []bool{true, false},
		StruFoo: struFoo,
	}

	struBarfoo := StruBarfoo{
		SS:   []string{"Foo", "Bar"},
		II:   []int{256, 512},
		BB:   []bool{true, false},
		SFoo: struFoo,
	}

	fmt.Println("If all Types X, Y, Z have the same capability of F1(), F2() etc., then they share the same Interface I := {X | exists F1(), F2() etc such that X.F1(), X.F2() etc. is valid}")

	var i I
	i = struFoo
	i.F1()
	i = struFooBar
	i.F1()
	i = struBarfoo
	i.F1()
}
