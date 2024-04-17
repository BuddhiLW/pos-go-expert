package fundacao

import "fmt"

func (s *StruFoo) FM1() {
	s.S = "passed through FM1"
	fmt.Println("Passing through FM1, we have:", s.S)
}

func (s StruFoo) FnM2() {
	s.S = "passed through FM2"
	fmt.Println("Passing through FnM2, we have:", s.S)
}

func pointers() {
	i := StruFoo{
		S: "Foo",
		I: 256,
		B: true,
	}

	fmt.Println("Pointers:")
	fmt.Println("i value is:", i)
	p := &i
	p.FM1()
	fmt.Println("p is a pointer for i:", p)
	fmt.Println("after p.FM1(), we have: i.S =", i.S)
	i.FnM2()
	fmt.Println("after i.FnM2(), we have: i.S =", i.S)
}
