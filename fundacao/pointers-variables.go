package fundacao

import "fmt"

func (s *StruFoo) FM1() {
	s.S = "passed through FM1"
	fmt.Println("While passing through FM1, we have: i.S =", s.S)
}

func (s StruFoo) FnM2() {
	s.S = "passed through FM2"
	fmt.Println("While passing through FnM2, we have: i.S =", s.S)
}

func Pointers() {
	i := StruFoo{
		S: "Foo",
		I: 256,
		B: true,
	}

	a := 256
	fmt.Printf("a value is: %v (a). It's adress is: %v (&a) \n", a, &a)

	fmt.Println("i value is:", i)
	// A function which Mutates
	i.FM1()
	fmt.Println("after i.FM1(), we have: i.S =", i.S)

	// A function which does not Mutate
	i.FnM2()
	fmt.Println("after i.FnM2(), we have: i.S =", i.S)

}
