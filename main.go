package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	tmpls "github.com/julvo/go-web-skeleton/templates"
)

var templates tmpls.Templates

func GetIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	templates.Execute("index.html", w, nil)
}

func main() {
	var err error
	templates, err = tmpls.New("./templates")
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.GET("/", GetIndex)
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
