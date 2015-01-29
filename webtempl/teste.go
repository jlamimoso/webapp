package main

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"html/template"
	"log"
	"net/http"
	"time"
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

/*
type hr struct {
	*httprouter.Router
}

func NewRouter() *hr {
	return &hr{httprouter.New()}
}
*/

type hr *httprouter.Router

func (r *hr) get(path string, h http.Handler) {
	r.GET(path, wrapHandler(h))
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		context.Set(r, "params", p)
		h.ServeHTTP(w, r)
	}
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("log -> [%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are on the about page 1.")
}

func main() {
	commonHandler := alice.New(loggingHandler)
	//r := NewRouter()
	r := hr(httprouter.New())
	r.get("/about", commonHandler.ThenFunc(aboutHandler))
	/*
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
	*/
	log.Fatal(http.ListenAndServe(":8080", r))
}
