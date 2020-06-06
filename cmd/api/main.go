package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/json-iterator/go"

	"github.com/gosom/context-spell-correct/pkg/spellcorrect"
)

type Result struct {
	Query       string
	Suggestions []spellcorrect.Suggestion
}

func getSuggestions(sc *spellcorrect.SpellCorrector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["query"]

		if !ok || len(keys[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		suggestions := sc.SpellCorrect(keys[0])

		cnt := 10
		if len(suggestions) < 10 {
			cnt = len(suggestions)
		}
		result := Result{
			Query:       keys[0],
			Suggestions: suggestions[:cnt],
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsoniter.NewEncoder(w).Encode(result)
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
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

	log.Fatal(http.ListenAndServe(":10000", logRequest(http.DefaultServeMux)))

}
