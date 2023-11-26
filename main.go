package main

import (
	fonctions "HangmanWeb/fonctions"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strings"
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
		Life       int
		TriedWords []string
		Setts      Settings
		Status     string
	}

	Actual := Settings{"easy", "french"}

	FileName := Actual.Difficulty + ".txt"
	FilePath := "texts/" + Actual.Language + "/" + FileName

	Game := GameData{FilePath, []string{}, "", "", 10, []string{}, Actual, "no"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(FilePath)
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
		FilePath = "texts/" + Actual.Language + "/" + FileName

		Game.Setts = Actual

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		if Game.Word == "" {
			Game.Word = fonctions.GetWord(fonctions.GetWords(FilePath))
			firstDisplay := fonctions.GetFirstDisplay(Game.Word)
			Game.Display = firstDisplay[0]
			Game.Letters = append(Game.Letters, firstDisplay[1])
		}
		Game.Status = "running"

		temp.ExecuteTemplate(w, "game", Game)
	})

	http.HandleFunc("/game/treatment", func(w http.ResponseWriter, r *http.Request) {
		try := r.FormValue("try")
		checkvalue, _ := regexp.MatchString("^[a-zA-Z-]{1-64}$", try)
		if !checkvalue {
			http.Redirect(w, r, "/game", http.StatusSeeOther)
		}
		try = strings.ToLower(try)
		if len(try) == 1 {
			Game.Letters = append(Game.Letters, try)
			if fonctions.VerifyLetter(Game.Word, try) {
				Game.Display = fonctions.Display(Game.Word, Game.Letters)
				if Game.Display == Game.Word {
					Game.Status = "won"
					http.Redirect(w, r, "/game/win", http.StatusSeeOther)
				}
			} else {
				Game.Life--
				if Game.Life <= 0 {
					Game.Status = "lost"
					http.Redirect(w, r, "/game/lose", http.StatusSeeOther)
				} else {
					http.Redirect(w, r, "/game", http.StatusSeeOther)
				}
			}
		} else if len(try) == len(Game.Word) {
			if try == Game.Word {
				Game.Status = "won"
				http.Redirect(w, r, "/game/win", http.StatusSeeOther)
			} else {
				Game.Life -= 2
				Game.TriedWords = append(Game.TriedWords, try)
				if Game.Life <= 0 {
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
