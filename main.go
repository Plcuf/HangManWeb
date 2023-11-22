package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func main() {
	temp, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		fmt.Println("Erreur > ", err)
		return
	}

	type Settings struct {
		Difficulty int
		Language   string
	}

	Actual := Settings{0, "English"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := Actual

		temp.ExecuteTemplate(w, "index", data)
	})

	http.HandleFunc("/rules", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "rules", nil)
	})

	RootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(RootDoc + "/web/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileserver))
	http.ListenAndServe("localhost:6969", nil)
}
