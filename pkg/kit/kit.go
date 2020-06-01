package kit

import (
	"bufio"
	//"fmt"
	"io"
	"strings"
	"sync"

	"gopkg.in/neurosnap/sentences.v1"
)

func Tokenize(in io.Reader) ([]string, error) {
	var tokens []string
	puncStrings := sentences.NewPunctStrings()
	tokenizer := sentences.NewWordTokenizer(puncStrings)
	linesc, errc := ReadIntoChan(in)
	r := strings.NewReplacer("(", " ", ")", " ")
	for line := range linesc {
		line = r.Replace(line)
		_tokens := tokenizer.Tokenize(line, false)
		for i := range _tokens {
			tokens = append(tokens, _tokens[i].Tok)
		}
	}
	err := <-errc
	return tokens, err

}

func ReadIntoChan(in io.Reader) (<-chan string, <-chan error) {
	outc := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(outc)
		defer close(errc)
		scanner := bufio.NewScanner(in)
		i := 0
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			if line != "" {
				outc <- scanner.Text()
			}
			i++
			if i == 5 {
				break
			}
		}
		if err := scanner.Err(); err != nil {
			errc <- err
			return
		}
	}()
	return outc, errc
}

func CombineErrors(cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error, len(cs))
	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func WaitErrors(errs ...<-chan error) error {
	for err := range CombineErrors(errs...) {
		if err != nil {
			return err
		}
	}
	return nil
}
