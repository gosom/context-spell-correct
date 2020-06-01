package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gosom/context-spell-correct/pkg/spellcorrect"
)

func main() {
	file, err := os.Open("datasets/sentences.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file2, err := os.Open("datasets/de-100k.txt")
	if err != nil {
		panic(err)
	}
	defer file2.Close()

	tokenizer := spellcorrect.NewSimpleTokenizer()
	freq := spellcorrect.NewFrequencies(2, 5)
	sc := spellcorrect.NewSpellCorrector(tokenizer, freq)
	t0 := time.Now()
	sc.Train(file, file2)
	t1 := time.Now()
	fmt.Printf("time to train %s\n", t1.Sub(t0))

	s := []string{
		"schloser in nurnberg",
		"schloser in karlruhe",
		"Rechtsawalte in Koln",
		"susi in leipzg",
		"Textilpflege in Bielefled",
		"textilpflege in frankfrt am main",
	}

	var search []string
	for i := 0; i < 10000; i++ {
		search = append(search, s...)
	}

	t2 := time.Now()
	for i := range search[:6000] {
		_ = sc.SpellCorrect(search[i])
	}
	t3 := time.Now()
	fmt.Printf("time to search %d items %s\n", len(search), t3.Sub(t2))

}
