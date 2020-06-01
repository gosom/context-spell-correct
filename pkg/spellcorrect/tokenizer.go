package spellcorrect

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

type SimpleTokenizer struct {
}

func NewSimpleTokenizer() *SimpleTokenizer {
	ans := SimpleTokenizer{}
	return &ans
}

func (o *SimpleTokenizer) Tokens(in io.Reader) ([]string, error) {
	// TODO do something more advanced here like use custom split function
	// that tokenizes properly
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanWords)
	var ans []string
	for scanner.Scan() {
		s := scanner.Text()
		s = strings.TrimRightFunc(s, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
		ans = append(ans, strings.ToLower(s))
	}

	err := scanner.Err()

	return ans, err
}
