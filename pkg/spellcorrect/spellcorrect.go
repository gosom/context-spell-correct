package spellcorrect

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/eskriett/spell"
)

type Suggestion struct {
	score  float64
	tokens []string
}

func (o *Suggestion) GetScore() float64 {
	return o.score
}

func (o *Suggestion) GetTokens() []string {
	return o.tokens
}

type FrequencyContainer interface {
	Load(tokens []string) error
	Get(tokens []string) float64
}

type Tokenizer interface {
	Tokens(in io.Reader) ([]string, error)
}

type SpellCorrector struct {
	tokenizer   Tokenizer
	frequencies FrequencyContainer
	spell       *spell.Spell
}

func NewSpellCorrector(tokenizer Tokenizer, frequencies FrequencyContainer) *SpellCorrector {
	ans := SpellCorrector{
		tokenizer:   tokenizer,
		frequencies: frequencies,
		spell:       spell.New(),
	}
	return &ans
}

func (o *SpellCorrector) Train(in io.Reader, in2 io.Reader) error {
	t0 := time.Now()
	tokens, err := o.tokenizer.Tokens(in)
	if err != nil {
		return err
	}
	t1 := time.Now()
	fmt.Println("time load tokens", t1.Sub(t0), len(tokens))
	o.frequencies.Load(tokens)
	t2 := time.Now()
	fmt.Println("time to load frequencies", t2.Sub(t1))

	scanner := bufio.NewScanner(in2)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		freq, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			panic(err)
		}
		o.spell.AddEntry(spell.Entry{
			Frequency: freq,
			Word:      parts[0],
		})

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return nil
}

type candidate struct {
	tokens        []string
	misspelledPos int
}

func candidatesEqual(a, b candidate) bool {
	if len(a.tokens) != len(b.tokens) {
		return false
	}
	for i := range a.tokens {
		if a.tokens[i] != b.tokens[i] {
			return false
		}
	}
	return true
}

func (o *SpellCorrector) getSuggestionCandidates(tokens []string, sugMap map[int]spell.SuggestionList, n int) []candidate {
	// TODO Improve me
	var candidates []candidate
	if n == 0 {
		return candidates
	}
	for i := 0; i < len(tokens); i++ {
		suggestions := sugMap[i]
		for j := range suggestions {
			b := append(tokens[:0:0], tokens...)
			b[i] = suggestions[j].Word
			cand := candidate{tokens: b}
			candidates = append(candidates, cand)
			candidates = append(candidates, o.getSuggestionCandidates(b, sugMap, n-1)...)
		}
	}
	return candidates
}

func (o *SpellCorrector) SpellCorrect(s string) []Suggestion {
	// TODO this code is bad -> improve
	tokens, _ := o.tokenizer.Tokens(strings.NewReader(s))
	allSuggestions := make(map[int]spell.SuggestionList)
	for i := range tokens {
		suggestions, _ := o.spell.Lookup(tokens[i], spell.SuggestionLevel(spell.LevelClosest))
		allSuggestions[i] = suggestions
	}
	dups := o.getSuggestionCandidates(tokens, allSuggestions, len(tokens))
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

	suggestions := o.rank(candidates)
	return suggestions
}

func (o *SpellCorrector) score(cand candidate) float64 {
	weights := []float64{5, 20, 2}
	score := 0.0
	for i := 1; i < 4; i++ {
		grams := TokenNgrams(cand.tokens, i)
		sum1 := 0.
		for i := range grams {
			sum1 += o.frequencies.Get(grams[i])
		}
		score += weights[i-1] * sum1
	}

	return score
}

func (o *SpellCorrector) rank(candidates []candidate) []Suggestion {
	var ans []Suggestion
	ans = make([]Suggestion, 0, len(candidates))

	for i := range candidates {
		score := o.score(candidates[i])
		ans = append(ans, Suggestion{score, candidates[i].tokens})
	}
	sort.SliceStable(ans, func(i, j int) bool {
		return ans[i].score > ans[j].score
	})
	if len(ans) > 10 {
		return ans[:10]
	}
	return ans
}
