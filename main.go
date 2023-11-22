package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func main() {
	temp, err := template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println("Erreur > ", err)
		return
	}

	type Data struct {
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := "coucou"

		temp.ExecuteTemplate(w, "index", data)
	})

	RootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(RootDoc + "/web/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileserver))
	http.ListenAndServe("localhost: 8080", nil)
}
