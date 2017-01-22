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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./tmpl/shared/head.tmpl", "./tmpl/home.tmpl"))
	t.ExecuteTemplate(w, "home", nil)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./tmpl/shared/head.tmpl", "./tmpl/404.tmpl"))
	t.ExecuteTemplate(w, "404", nil)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	post := Post{}
	post.Title = r.URL.Path[len("/posts/"):]

	input, err := ioutil.ReadFile("./posts/" + post.Title + ".md")
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusNotFound)
		return
	}

	output := blackfriday.MarkdownCommon(input)
	post.Body = template.HTML(output)

	t := template.Must(template.ParseFiles("./tmpl/shared/head.tmpl", "./tmpl/post.tmpl"))
	t.ExecuteTemplate(w, "post", post)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/posts/", PostHandler)
	http.HandleFunc("/404/", NotFoundHandler)

	http.ListenAndServe(":8080", nil)
}
