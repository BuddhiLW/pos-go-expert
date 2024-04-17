package fundacao

import "fmt"

// Custom Type StruFoo
type StruFoo struct {
	S string
	I int
	B bool
}

// Composed custom type StruFooBar
type StruFooBar struct {
	SS []string
	II []int
	BB []bool
	StruFoo
}

// Custom type, StruBarfoo, relying on another Custom Type, StruFoo
type StruBarfoo struct {
	SS   []string
	II   []int
	BB   []bool
	SFoo StruFoo
}

func Structs() {
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
	fmt.Printf("struFoo: %v %T \n", struFoo, struFoo)
	fmt.Printf("struFooBar: %v %T \n", struFooBar, struFooBar)
	fmt.Printf("struBarfoo: %v %T \n\n", struBarfoo, struBarfoo)

	fmt.Println("Struct instances of StruFooBar and StruBarfoo will look identical, but aren't.")
	fmt.Println("Composition vs Standard type (implies direct access vs not direct access):")
	fmt.Printf("struFooBar.S (%v) == struBarfoo.SFoo.S (%v)? %v", struFooBar.S, struBarfoo.SFoo.S, struFooBar.S == struBarfoo.SFoo.S)
}
