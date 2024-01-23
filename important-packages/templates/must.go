package templates

import (
	"fmt"
	"os"
	"text/template"
)

func Must() {
	fmt.Println("Templating using template.Must()")
	course := Course{"Go Expert", 1080}
	t := template.Must(template.New("CourseTemplate").Parse("Name: {{.Name}} - Credit (h): {{.Credit}}\n"))
	err := t.Execute(os.Stdout, course)
	if err != nil {
		panic(err)
	}
}
