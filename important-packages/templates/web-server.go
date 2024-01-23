package templates

import (
	"net/http"
	"text/template"
)

func TemplateWebServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("courses.html").ParseFiles(Wd() + "/templates/courses.html"))
		err := t.Execute(w, Courses{
			{"Go Expert", 360},
			{"Go Expert II", 720},
			{"Go Expert III", 1080},
		})

		if err != nil {
			panic(err)
		}
	})
	http.ListenAndServe(":8080", nil)
}
