package spellcorrect

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/eskriett/spell"
)

type Suggestion struct {
	score  float64
	Tokens []string
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

func hashTokens(tokens []string) uint64 {
	h := fnv.New64a()
	for i := range tokens {
		h.Write([]byte(tokens[i]))
	}
	return h.Sum64()
}

func product(a []string, b []string) []string {
	var items []string
	for i := range a {
		for j := range b {
			items = append(items, a[i]+" "+b[j])
		}
	}
	return items
}

func combos(in [][]string) []string {
	tmpP := in[len(in)-1]
	for i := len(in) - 2; i >= 0; i-- {
		tmpP = product(in[i], tmpP)
	}
	return tmpP
}

func (o *SpellCorrector) lookupTokens(tokens []string) [][]string {
	var allSuggestions [][]string
	for i := range tokens {
		allSuggestions = append(allSuggestions, nil)
		allSuggestions[i] = append(allSuggestions[i], tokens[i])
		suggestions, _ := o.spell.Lookup(tokens[i], spell.SuggestionLevel(spell.LevelClosest))
		for j := 0; j < len(suggestions) && j < 10; j++ {
			allSuggestions[i] = append(allSuggestions[i], suggestions[j].Word)
		}
	}
	return allSuggestions
}

func (o *SpellCorrector) getSuggestionCandidates(allSuggestions [][]string) []Suggestion {
	suggestionStrings := combos(allSuggestions)
	seen := make(map[uint64]struct{}, len(suggestionStrings))
	suggestions := make([]Suggestion, 0, len(suggestionStrings))
	for i := range suggestionStrings {
		sugTokens := strings.Split(suggestionStrings[i], " ")
		h := hashTokens(sugTokens)
		if _, ok := seen[h]; !ok {
			seen[h] = struct{}{}
			suggestions = append(suggestions,
				Suggestion{
					score:  o.score(sugTokens),
					Tokens: sugTokens,
				})
		}
	}
	sort.SliceStable(suggestions, func(i, j int) bool {
		return suggestions[i].score > suggestions[j].score
	})
	return suggestions
}

func (o *SpellCorrector) SpellCorrect(s string) []Suggestion {
	tokens, _ := o.tokenizer.Tokens(strings.NewReader(s))
	allSuggestions := o.lookupTokens(tokens)
	return o.getSuggestionCandidates(allSuggestions)
}

func (o *SpellCorrector) score(tokens []string) float64 {
	weights := []float64{100, 15, 5}
	score := 0.0
	for i := 1; i < 4; i++ {
		grams := TokenNgrams(tokens, i)
		sum1 := 0.
		for i := range grams {
			sum1 += o.frequencies.Get(grams[i])
		}
		score += weights[i-1] * sum1
	}
	return score
}
