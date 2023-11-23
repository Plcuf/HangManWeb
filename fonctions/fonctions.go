package fonctions

import (
	"fmt"
	"math/rand"
	"os"
)

func GetWords(name string) []string {
	file, err := os.OpenFile(name, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Erreur > ", err)
	}
	defer file.Close()

	data, err := os.ReadFile(file.Name())
	if err != nil {
		fmt.Println("Erreur > ", err)
	}

	slice := []string{}
	word := ""

	for _, c := range data {
		c := string(c)
		if c == "\n" {
			slice = append(slice, word)
			word = ""
		} else {
			word = word + c
		}
	}
	return slice
}

func GetWord(s []string) string {
	brokenword := s[rand.Intn(len(s))]
	word := ""
	for i := 0; i <= len(brokenword)-1; i++ {
		word += string(brokenword[i])
	}
	return word
}

func GetFirstDisplay(word string) string {
	displayedLetter := word[rand.Int31n(int32(len(word)-1))]
	display := ""
	for i := 0; i < len(word); i++ {
		if word[i] == displayedLetter {
			display += string(word[i])
		} else {
			display += "_"
		}
	}
	return display
}

func VerifyLetter(word string, letter string) bool {
	for i := 0; i < len(word); i++ {
		if letter[0] == word[i] {
			return true
		}
	}
	return false
}

func Display(word string, letter byte, displayed string) string {
	display := ""
	for i := 0; i < len(word); i++ {
		if word[i] == letter {
			display += string(letter)
		} else {
			display += displayed[i:]
		}
	}
	return display
}