package main

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/russross/blackfriday"
)

type Post struct {
	Title string
	Body  template.HTML
}

var t = template.New("views")

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "home", nil)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "404", nil)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	post := Post{}
	post.Title = r.URL.Path[len("/posts/"):]

	input, err := ioutil.ReadFile("./posts/" + post.Title + ".md")
	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	output := blackfriday.MarkdownCommon(input)
	post.Body = template.HTML(output)

	t.ExecuteTemplate(w, "post", post)
}

func LoadViews() {
	t = template.Must(t.ParseGlob("./tmpl/*.tmpl"))
	t = template.Must(t.ParseGlob("./tmpl/shared/*.tmpl"))
}

func main() {
	LoadViews()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/posts/", PostHandler)
	http.HandleFunc("/404/", NotFoundHandler)

	http.ListenAndServe(":8080", nil)
}
