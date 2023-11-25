package fonctions

func VerifyLetter(word string, letter string) bool {
	isGood := false

	for i := 0; i < len(word); i++ {
		if letter == string(word[i]) {
			isGood = true
		}
	}

	return isGood
}

func VerifyWord(word string, tried string) bool {
	return word == tried
}
