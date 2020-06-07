package webhandlers

import (
	"log"
	"net/http"

	"github.com/json-iterator/go"

	"github.com/gosom/context-spell-correct/pkg/spellcorrect"
)

type Result struct {
	Query       string
	Suggestions []spellcorrect.Suggestion
}

func GetSuggestions(sc *spellcorrect.SpellCorrector) http.HandlerFunc {
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

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
