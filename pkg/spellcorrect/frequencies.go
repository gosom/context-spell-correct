package spellcorrect

import (
	"fmt"
	"time"

	"github.com/segmentio/fasthash/fnv1a"
)

type ngram []uint64

type Frequencies struct {
	minWord      int
	minFreq      int
	uniGramProbs map[uint64]float64
	trie         *wordTrie
}

func NewFrequencies(minWord, minFreq int) *Frequencies {
	ans := Frequencies{
		minWord:      minWord,
		minFreq:      minFreq,
		uniGramProbs: make(map[uint64]float64),
		trie:         newWordTrie(0),
	}
	return &ans
}

func (o *Frequencies) Load(tokens []string) error {
	o.trie = newWordTrie(len(tokens))
	t1 := time.Now()
	hashes := make([]uint64, len(tokens), len(tokens))
	bl := make(map[uint64]bool)
	unigrams := make(map[uint64]int)
	for i := range tokens {
		hashes[i] = hashString(tokens[i])
		unigrams[hashes[i]]++
		if len(tokens[i]) < o.minWord {
			bl[hashes[i]] = true
		}
	}

	for k, v := range unigrams {
		if v < o.minFreq {
			bl[k] = true
		} else {
			o.uniGramProbs[k] = float64(v) / float64(len(tokens))
		}
	}

	t2 := time.Now()
	fmt.Println("time to hash and map", t2.Sub(t1))

	for i := 1; i < 4; i++ {
		grams := ngrams(hashes, i)
		for _ngram := range grams {
			add := true
			for j := range _ngram {
				if bl[_ngram[j]] {
					add = false
					break
				}
			}
			if add {
				o.trie.put(_ngram)
			}
		}
	}
	t3 := time.Now()
	fmt.Println("Time to add to trie", t3.Sub(t2))

	return nil
}

func (o *Frequencies) Get(tokens []string) float64 {
	hashes := make([]uint64, len(tokens), len(tokens))
	for i := range tokens {
		hashes[i] = hashString(tokens[i])
	}
	if len(hashes) == 1 {
		return o.uniGramProbs[hashes[0]]
	}
	node := o.trie.search(hashes)
	if node == nil {
		return 0.0
	}
	return node.prob
}

type node struct {
	freq     int
	prob     float64
	children map[uint64]*node
}

func newNode(freq int) *node {
	n := node{
		freq:     freq,
		children: make(map[uint64]*node),
	}
	return &n
}

type wordTrie struct {
	root *node
}

func newWordTrie(lenTokens int) *wordTrie {
	trie := wordTrie{
		root: newNode(lenTokens),
	}
	return &trie
}

//The assumption that we first add the 1gram then the 2gram etc is made
func (o *wordTrie) put(key ngram) {
	current := o.root
	for i := 0; i < len(key); i++ {
		if i == len(key)-1 {
			node, ok := current.children[key[i]]
			if ok {
				node.freq++
			} else {
				node = newNode(1)
				current.children[key[i]] = node
			}
			node.prob = float64(node.freq) / float64(current.freq)
		} else {
			current = current.children[key[i]]
		}
	}
}

func (o *wordTrie) search(key ngram) *node {
	tmp := o.root
	for i := 0; i < len(key); i++ {
		if next, ok := tmp.children[key[i]]; ok {
			tmp = next
		} else {
			return nil
		}
	}
	return tmp
}

func hashString(s string) uint64 {
	return fnv1a.HashString64(s)
}

func TokenNgrams(words []string, size int) [][]string {
	var out [][]string
	for i := 0; i+size <= len(words); i++ {
		out = append(out, words[i:i+size])
	}
	return out
}

func ngrams(words []uint64, size int) <-chan ngram {
	out := make(chan ngram)
	go func() {
		defer close(out)
		for i := 0; i+size <= len(words); i++ {
			out <- words[i : i+size]
		}
	}()
	return out
}
