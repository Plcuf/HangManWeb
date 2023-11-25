package fonctions

import "math/rand"

func GetFirstDisplay(word string) string {
	discovered := rand.Intn(len(word))
	displayed := ""

	for i := 0; i < len(word); i++ {
		if i == discovered {
			displayed += string(word[i])
		} else {
			displayed += "_"
		}
	}

	return displayed
}

func Display(word string, discovered []string) string {
	displayed := ""

	for i := 0; i < len(word); i++ {
		letterDisplayed := false
		for j := 0; j < len(discovered); j++ {
			if string(word[i]) == discovered[j] {
				letterDisplayed = true
			}
		}
		if letterDisplayed {
			displayed += string(word[i])
		} else {
			displayed += "_"
		}
	}

	return displayed
}
