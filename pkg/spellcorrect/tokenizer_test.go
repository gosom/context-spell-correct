package spellcorrect

import (
	"strings"
	"testing"
)

func TestSimpleTokenizerTokens(t *testing.T) {
	traindata := `Als das Haus in sich zusammenbrach, wurden viele der Opfer 
				von herabstürzenden Teilen getroffen`

	expected := []string{
		"als", "das", "haus", "in", "sich", "zusammenbrach", "wurden", "viele",
		"der", "opfer", "von", "herabstürzenden", "teilen", "getroffen",
	}

	r := strings.NewReader(traindata)
	tokenizer := NewSimpleTokenizer()
	tokens, err := tokenizer.Tokens(r)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if len(tokens) != len(expected) {
		t.Errorf("tokens differ (len)")
		return
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != tokens[i] {
			t.Errorf("token (%s) in position %d differ", expected[i], i)
			return
		}
	}

}
