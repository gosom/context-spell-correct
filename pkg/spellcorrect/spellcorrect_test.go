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
	sc := NewSpellCorrector(tokenizer, freq, []float64{100, 15, 5})
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

	var allSuggestions [][]string
	for i := range tokens {
		allSuggestions = append(allSuggestions, nil)
		allSuggestions[i] = append(allSuggestions[i], tokens[i])
		suggestions, _ := sugMap[i]
		for j := 0; j < len(suggestions) && j < 10; j++ {
			allSuggestions[i] = append(allSuggestions[i], suggestions[j].Word)
		}

	}

	expected := [][]string{

		[]string{"1", "2", "3"},
		[]string{"a", "2", "3"},
		[]string{"aa", "2", "3"},
		[]string{"1", "b", "3"},
		[]string{"a", "b", "3"},
		[]string{"aa", "b", "3"},
	}

	sc := getSpellCorrector()

	candidates := sc.getSuggestionCandidates(allSuggestions)

	if len(candidates) != len(expected) {
		t.Errorf("invalid length")
		return
	}

	expect := make(map[uint64]bool)
	for i := range expected {
		expect[hashTokens(expected[i])] = true
	}
	for i := range candidates {
		if !expect[hashTokens(candidates[i].Tokens)] {
			t.Errorf("%v not in expected", candidates[i].Tokens)
			return
		}
	}

}
