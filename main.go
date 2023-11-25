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

	RootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(RootDoc + "/web/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileserver))

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
		Setts      Settings
		Status     string
	}

	Actual := Settings{"easy", "french"}

	FileName := Actual.Difficulty + ".txt"
	FilePath := "/web/texts/" + Actual.Language + "/" + FileName

	Game := GameData{FilePath, []string{}, "", "", 0, 10, []string{}, Actual, "no"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "index", nil)
	})

	http.HandleFunc("/rules", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "rules", nil)
	})

	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "settings", Actual)
	})

	http.HandleFunc("/settings/treatment", func(w http.ResponseWriter, r *http.Request) {
		Actual = Settings{r.FormValue("difficulty"), r.FormValue("language")}

		FileName = Actual.Difficulty + ".txt"
		FilePath = "/web/texts/" + Actual.Language + "/" + FileName

		Game.Setts = Actual

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		if Game.Word == "" {
			Game.Word = fonctions.GetWord(fonctions.GetWords(Game.File))
			Game.Display = fonctions.GetFirstDisplay(Game.Word)
		}
		Game.Status = "running"

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
					Game.Status = "won"
					http.Redirect(w, r, "/game/win", http.StatusSeeOther)
				}
			} else {
				Game.Life--
				if Game.Life == 0 {
					Game.Status = "lost"
					http.Redirect(w, r, "/game/lose", http.StatusSeeOther)
				} else {
					http.Redirect(w, r, "/game", http.StatusSeeOther)
				}
			}
		} else if len(triedWord) == len(Game.Word) {
			if triedWord == Game.Word {
				Game.Status = "won"
				http.Redirect(w, r, "/game/win", http.StatusSeeOther)
			} else {
				Game.Life--
				Game.TriedWords = append(Game.TriedWords, triedWord)
				if Game.Life == 0 {
					Game.Status = "lost"
					http.Redirect(w, r, "/game/lose", http.StatusSeeOther)
				} else {
					http.Redirect(w, r, "/game", http.StatusSeeOther)
				}
			}
		}
	})

	http.HandleFunc("/game/win", func(w http.ResponseWriter, r *http.Request) {
		if Game.Status != "won" {
			http.Redirect(w, r, "/game", http.StatusSeeOther)
		} else {
			Game.Life = 10
			Game.Word = ""
			temp.ExecuteTemplate(w, "win", Game)
		}
	})

	http.HandleFunc("/game/lose", func(w http.ResponseWriter, r *http.Request) {
		if Game.Status != "lost" {
			http.Redirect(w, r, "/game", http.StatusSeeOther)
		} else {
			Game.Life = 10
			Game.Word = ""
			temp.ExecuteTemplate(w, "lose", Game)
		}
	})

	http.ListenAndServe("localhost:6969", nil)
}
