package templates

import (
	"fmt"
	"os"
	"text/template"
)

type Course struct {
	Name   string
	Credit int
}

func First() {
	fmt.Println("First template; prints the =Name= and =Credit= information of a course.")
	course := Course{"Go Expert", 1080}
	tmp := template.New("CourseTemplate")
	tmp, _ = tmp.Parse("Name: {{.Name}} - Credit (h): {{.Credit}}\n")
	err := tmp.Execute(os.Stdout, course)
	if err != nil {
		panic(err)
	}
}
