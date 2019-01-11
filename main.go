package main

import (
    "net/http"
    "html/template"
    "log"

    "path/filepath"

    "github.com/julienschmidt/httprouter"
)

var templates map[string]*template.Template

func GetIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    templates["index.html"].Execute(w, map[string]interface{}{
    })
}


func main() {
    var err error
    templates, err = LoadTemplates("./templates")
    if err != nil {
        log.Fatal(err)
    }

    router := httprouter.New()
    router.GET("/", GetIndex)

    router.ServeFiles("/static/*filepath", http.Dir("static"))

    log.Fatal(http.ListenAndServe(":8080", router))
}

func LoadTemplates(dir string) (map[string]*template.Template, error) {
    partials, err := filepath.Glob(filepath.Join(dir, "partials", "*.html"))
    if err != nil {
        return nil, err
    }
    pages, err := filepath.Glob(filepath.Join(dir, "pages", "*.html"))
    if err != nil {
        return nil, err
    }
    root := filepath.Join(dir, "root.html")

    templates := make(map[string]*template.Template)
    for _, page := range(pages) {
        relPath, _ := filepath.Rel(filepath.Join(dir, "pages"), page)
        templ, err := template.ParseFiles(append([]string{root, page}, partials...)...)
        if err != nil {
            return nil, err
        }
        templates[relPath] = templ
    }
    return templates, nil
}

