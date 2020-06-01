package spellcorrect

import (
	"fmt"
	"strings"
	"testing"

	"github.com/eskriett/spell"
)

func getSpellCorrector() *SpellCorrector {
	tokenizer := NewSimpleTokenizer()
	freq := NewFrequencies(0, 0)
	sc := NewSpellCorrector(tokenizer, freq)
	return sc
}

func TestTrain(t *testing.T) {
	trainwords := "golang 100\ngoland 1"
	traindata := `golang python C erlang golang java java golang goland`
	r := strings.NewReader(traindata)
	r1 := strings.NewReader(trainwords)

	sc := getSpellCorrector()
	if err := sc.Train(r, r1); err != nil {
		t.Errorf(err.Error())
		return
	}
	if prob := sc.frequencies.Get([]string{"golang"}); prob != 0.5 {
		t.Errorf("invalid prob %f", prob)
		return
	}

	suggestions, _ := sc.spell.Lookup("gola", spell.SuggestionLevel(spell.LevelAll))
	fmt.Println(suggestions)
}

func TestSpellCorrect(t *testing.T) {
	trainwords := "golang 100\ngoland 1"
	traindata := `golang python C erlang golang java java golang goland`
	r := strings.NewReader(traindata)
	r1 := strings.NewReader(trainwords)

	sc := getSpellCorrector()
	if err := sc.Train(r, r1); err != nil {
		t.Errorf(err.Error())
		return
	}

	s1 := "restaurant in Bonn"

	suggestions := sc.SpellCorrect(s1)
	fmt.Println(suggestions)
}

func TestGetSuggestionCandidates(t *testing.T) {
	tokens := []string{"1", "2", "3"}

	sugMap := map[int]spell.SuggestionList{
		0: spell.SuggestionList{
			spell.Suggestion{Distance: 1, Entry: spell.Entry{Word: "a"}},
			spell.Suggestion{Distance: 2, Entry: spell.Entry{Word: "aa"}},
		},
		1: spell.SuggestionList{
			spell.Suggestion{Distance: 1, Entry: spell.Entry{Word: "b"}},
		},
		2: spell.SuggestionList{},
	}

	expected := []candidate{
		// 0
		candidate{[]string{"a", "2", "3"}, 0},
		candidate{[]string{"aa", "2", "3"}, 0},

		// 1
		candidate{[]string{"1", "b", "3"}, 0},
		candidate{[]string{"a", "b", "3"}, 0},
		candidate{[]string{"aa", "b", "3"}, 0},
	}

	sc := getSpellCorrector()

	dups := sc.getSuggestionCandidates(tokens, sugMap, len(tokens))

	var candidates []candidate
	for i := range dups {
		found := false
		for j := range candidates {
			if candidatesEqual(dups[i], candidates[j]) {
				found = true
				break
			}
		}
		if !found {
			candidates = append(candidates, dups[i])
		}
	}
	fmt.Println(candidates)

	if len(candidates) != len(expected) {
		t.Errorf("invalid length")
		return
	}

	for i := range expected {
		e := expected[i].tokens
		a := candidates[i].tokens
		if len(a) != len(e) {
			t.Errorf("invalid len of tokens")
			return
		}
		for j := range e {
			if e[j] != a[j] {
				t.Errorf("Token at %d (%s vs %s) differ", j, e[j], a[j])
				return
			}
		}
	}

}
