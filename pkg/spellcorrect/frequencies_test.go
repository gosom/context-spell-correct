package spellcorrect

import (
	"fmt"
	"testing"
)

func TestFrequencies(t *testing.T) {
	tokens := []string{"I", "program", "go", "I", "code", "and", "I", "cook", "code"}
	freq := NewFrequencies(0, 0)
	if err := freq.Load(tokens); err != nil {
		t.Errorf(err.Error())
		return
	}

	prob := freq.Get([]string{"I"})
	fmt.Println("I", prob)
	prob2 := freq.Get([]string{"I", "code"})
	fmt.Println("Second", prob2)

	prob3 := freq.Get([]string{"I", "program", "go"})
	fmt.Println(prob3)

	/*
		if prob := freq.GetProbability([]string{"I"}); prob != 110 || prob > 0.34 && prob < 0.33 {
			t.Errorf("calculated invalid probability %f", prob)
			return
		}

		arithm := freq.Get([]string{"I", "program"})
		fmt.Println(arithm)
		par := freq.Get([]string{"I"})
		fmt.Println(par)

		if prob := freq.GetProbability([]string{"I", "program"}); prob != 0 {
			t.Errorf("calculated invalid probabiliyu %f", prob)
		}
	*/
}

/*
func TestWordTrie(t *testing.T) {
	trie := newWordTrie()

	words := []uint64{
		1, 2, 3, 4, 5, 6, 1, 2,
	}

	unigrams := ngrams(words, 1)
	for i := range unigrams {
		trie.put(unigrams[i])
	}

	s := ngram{uint64(2)}
	if freq := trie.search(s); freq != 2 {
		t.Errorf("error computing freq")
		return
	}
	if freq := trie.search(ngram{uint64(79)}); freq != 0 {
		t.Errorf("error computing freq")
		return
	}
	bigrams := ngrams(words, 2)
	for i := range bigrams {
		trie.put(bigrams[i])
	}

	if freq := trie.search(ngram{uint64(1)}); freq != 2 {
		t.Errorf("error computing freq")
		return
	}
	if freq := trie.search(ngram{uint64(1), uint64(2)}); freq != 2 {
		t.Errorf("error computing freq")
		return
	}
}
*/
