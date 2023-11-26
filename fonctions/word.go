package fonctions

import (
	"math/rand"
	"os"
	"strings"
)

func GetWords(name string) []string {
	file, err := os.OpenFile(name, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := os.ReadFile(file.Name())
	if err != nil {
		panic(err)
	}

	slice := []string{}
	word := ""

	for _, c := range data {
		c := string(c)
		if c == "\n" {
			word = strings.TrimSpace(word)
			slice = append(slice, word)
			word = ""
		} else {
			word = word + c
		}
	}
	return slice
}

func GetWord(s []string) string {
	word := s[rand.Intn(len(s))]
	return word
}
