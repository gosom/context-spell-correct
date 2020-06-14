package main

import (
	"compress/gzip"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gosom/context-spell-correct/internal/config"
	"github.com/gosom/context-spell-correct/internal/webhandlers"
	"github.com/gosom/context-spell-correct/pkg/spellcorrect"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	file, err := os.Open(cfg.SentencesPath)
	if err != nil {
		panic(err)
	}

	gz, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}

	file2, err := os.Open(cfg.DictPath)
	if err != nil {
		panic(err)
	}
	gz2, err := gzip.NewReader(file2)
	if err != nil {
		panic(err)
	}

	tokenizer := spellcorrect.NewSimpleTokenizer()

	freq := spellcorrect.NewFrequencies(cfg.MinWordLength, cfg.MinWordFreq)

	weights := []float64{cfg.UnigramWeight, cfg.BigramWeight, cfg.TrigramWeight}
	sc := spellcorrect.NewSpellCorrector(tokenizer, freq, weights)

	log.Printf("starting training...")
	t0 := time.Now()
	sc.Train(gz, gz2)
	t1 := time.Now()
	log.Printf("ready[%s]\n", t1.Sub(t0))
	file.Close()
	gz.Close()
	file2.Close()
	gz2.Close()

	http.HandleFunc("/", webhandlers.GetSuggestions(sc))

	log.Fatal(http.ListenAndServe(cfg.Addr, webhandlers.LogRequest(http.DefaultServeMux)))

}
