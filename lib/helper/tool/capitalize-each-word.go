package tool

import "strings"

func CapitalizeEachWord(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, " ")
}
