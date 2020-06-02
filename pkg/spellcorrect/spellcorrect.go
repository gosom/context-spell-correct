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
	tokens []string
}

func hashCandidate(item candidate) uint64 {
	h := fnv.New64a()
	for i := range item.tokens {
		h.Write([]byte(item.tokens[i]))
	}
	return h.Sum64()
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

func (o *SpellCorrector) generateCandidates(tokens []string, candidates map[uint64]candidate, n int,
	sugMap map[int]spell.SuggestionList) {
	if n == len(tokens) {
		return
	}

	for i := 0; i < len(tokens); i++ {
		suggestions := sugMap[i]
		for _, sug := range suggestions {
			b := append(tokens[:0:0], tokens...)
			b[i] = sug.Word
			cand := candidate{tokens: b}
			candHash := hashCandidate(cand)
			if _, ok := candidates[candHash]; !ok {
				candidates[candHash] = cand
			}
			o.generateCandidates(b, candidates, n+1, sugMap)
		}
	}
}

func (o *SpellCorrector) getSuggestionCandidates(tokens []string, sugMap map[int]spell.SuggestionList) []candidate {

	found := make(map[uint64]candidate)

	o.generateCandidates(tokens, found, 0, sugMap)

	candidates := make([]candidate, 0, len(found))
	for _, v := range found {
		candidates = append(candidates, v)
	}
	return candidates
}

func (o *SpellCorrector) SpellCorrect(s string) []Suggestion {
	// maybe something more efficient
	tokens, _ := o.tokenizer.Tokens(strings.NewReader(s))
	allSuggestions := make(map[int]spell.SuggestionList)
	for i := range tokens {
		suggestions, _ := o.spell.Lookup(tokens[i], spell.SuggestionLevel(spell.LevelClosest))
		allSuggestions[i] = suggestions
	}
	candidates := o.getSuggestionCandidates(tokens, allSuggestions)

	suggestions := o.rank(candidates)
	return suggestions
}

func (o *SpellCorrector) score(cand candidate) float64 {
	weights := []float64{10, 15, 5}
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
