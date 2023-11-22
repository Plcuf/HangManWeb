package main

import (
	"HangmanWeb/fonctions"
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
		Difficulty string
		Language   string
	}

	type GameData struct {
		File    string
		Letters []string
		Word    string
	}

	Actual := Settings{"easy", "English"}

	FileName := Actual.Difficulty + ".txt"
	FilePath := "/assets/texts/" + Actual.Language + "/" + FileName

	Game := GameData{FilePath, []string{}, ""}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := Actual

		temp.ExecuteTemplate(w, "index", data)
	})

	http.HandleFunc("/rules", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "rules", nil)
	})

	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "settings", nil)
	})

	http.HandleFunc("/settings/treatment", func(w http.ResponseWriter, r *http.Request) {
		Actual = Settings{r.FormValue("difficulty"), r.FormValue("language")}

		FileName = Actual.Difficulty + ".txt"
		FilePath = "/assets/texts/" + Actual.Language + "/" + FileName

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		Game.Word = fonctions.GetWord(fonctions.GetWords())

		temp.ExecuteTemplate(w, "game", Game)
	})

	http.HandleFunc("/game/treatment", func(w http.ResponseWriter, r *http.Request) {
		Game.Letters = append(Game.Letters, r.FormValue("letter"))

		http.Redirect(w, r, "/game", http.StatusSeeOther)
	})

	RootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(RootDoc + "/web/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileserver))
	http.ListenAndServe("localhost:6969", nil)
}
