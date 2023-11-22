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
		File       string
		Letters    []string
		Word       string
		Display    string
		Points     int
		Life       int
		TriedWords []string
	}

	Actual := Settings{"easy", "english"}

	FileName := Actual.Difficulty + ".txt"
	FilePath := "/assets/texts/" + Actual.Language + "/" + FileName

	Game := GameData{FilePath, []string{}, "", "", 0, 10, []string{}}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Game.Word = ""

		temp.ExecuteTemplate(w, "index", nil)
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
		if Game.Word == "" {
			Game.Word = fonctions.GetWord(fonctions.GetWords(FilePath))
			Game.Display = fonctions.GetFirstDisplay(Game.Word)
		}

		temp.ExecuteTemplate(w, "game", Game)
	})

	http.HandleFunc("/game/treatment", func(w http.ResponseWriter, r *http.Request) {
		playedLetter := r.FormValue("letter")
		triedWord := r.FormValue("word")
		if len(playedLetter) == 1 {
			Game.Letters = append(Game.Letters, playedLetter)
			if fonctions.VerifyLetter(Game.Word, playedLetter) {
				Game.Display = fonctions.Display(Game.Word, playedLetter[0], Game.Display)
				if Game.Display == Game.Word {
					http.Redirect(w, r, "/game/win", http.StatusSeeOther)
				}
			} else {
				Game.Life--
				if Game.Life == 0 {
					http.Redirect(w, r, "/game/lose", http.StatusSeeOther)
				} else {
					http.Redirect(w, r, "/game", http.StatusSeeOther)
				}
			}
		} else if len(triedWord) == len(Game.Word) {
			if triedWord == Game.Word {
				http.Redirect(w, r, "/game/win", http.StatusSeeOther)
			} else {
				Game.Life--
				Game.TriedWords = append(Game.TriedWords, triedWord)
				if Game.Life == 0 {
					http.Redirect(w, r, "/game/lose", http.StatusSeeOther)
				} else {
					http.Redirect(w, r, "/game", http.StatusSeeOther)
				}
			}
		}

	})

	http.HandleFunc("/game/Win", func(w http.ResponseWriter, r *http.Request) {
		Game.Life = 10
		temp.ExecuteTemplate(w, "win", Game)
	})

	http.HandleFunc("/game/lose", func(w http.ResponseWriter, r *http.Request) {
		Game.Life = 10
		temp.ExecuteTemplate(w, "lose", Game)
	})

	RootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(RootDoc + "/web/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileserver))
	http.ListenAndServe("localhost:6969", nil)
}
