package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "index.html", "templates/home.html"))

type tipox func(int) int

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//fmt.Fprint(w, "Welcome!\n")
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func css(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h := http.StripPrefix("/css/", http.FileServer(http.Dir("./css/")))
	fmt.Printf("passe no cs !!!!\n")
	h.ServeHTTP(w, r)
}

func js(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h := http.StripPrefix("/js/", http.FileServer(http.Dir("./js/")))
	fmt.Printf("passe no js !!!!\n")
	h.ServeHTTP(w, r)
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func mostrarX(x int) int {
	return x + 1
}

func (x tipox) teste(v int) int {
	return x(v)
}

type router struct {
	*httprouter.Router
}

func NewRouter() *router {
	return &router{httprouter.New()}
}

func main() {
	r := NewRouter()

	router := httprouter.New()
	y := tipox(mostrarX)
	fmt.Printf("valor do tipo x %d", y.teste(2))
	router.ServeFiles("/css/*filepath", http.Dir("./css/"))
	router.ServeFiles("/js/*filepath", http.Dir("./js/"))
	router.GET("/", Index)
	//router.GET("/css/", css)
	//router.GET("/js/", js)
	router.GET("/hello/:name", Hello)

	log.Fatal(http.ListenAndServe(":8080", router))
}
