package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gosom/context-spell-correct/pkg/spellcorrect"
)

type Suggestion struct {
	Score  float64
	Tokens []string
}

type Result struct {
	Query       string
	Suggestions []Suggestion
}

func getSuggestions(sc *spellcorrect.SpellCorrector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["query"]

		if !ok || len(keys[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		suggestions := sc.SpellCorrect(keys[0])

		result := Result{
			Query:       keys[0],
			Suggestions: make([]Suggestion, len(suggestions)),
		}
		for i := range suggestions {
			result.Suggestions[i] = Suggestion{
				Score:  suggestions[i].GetScore(),
				Tokens: suggestions[i].GetTokens(),
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

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

	http.HandleFunc("/", getSuggestions(sc))
	log.Fatal(http.ListenAndServe(":10000", nil))

}
