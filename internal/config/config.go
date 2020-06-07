package config

import (
	"fmt"
	"log"

	"github.com/qiangxue/go-env"
)

type Config struct {
	Addr          string  `env:"ADDR"`
	SentencesPath string  `env:"SENTENCES_PATH"`
	DictPath      string  `env:"DICT_PATH"`
	MinWordLength int     `env:"MIN_WORD_LENGTH"`
	MinWordFreq   int     `env:"MIN_WORD_FREQ"`
	UnigramWeight float64 `env:"UNIGRAM_WEIGHT"`
	BigramWeight  float64 `env:"BIGRAM_WEIGHT"`
	TrigramWeight float64 `env:"TRIGRAM_WEIGHT"`
}

func (o *Config) Validate() error {
	if o.Addr == "" {
		return fmt.Errorf("provide non empty SC_ADDR")
	}
	if o.SentencesPath == "" {
		return fmt.Errorf("provide non empty SC_SENTENCES_PATH")
	}
	if o.DictPath == "" {
		return fmt.Errorf("provide non empty SC_DICT_PATH")
	}
	if o.MinWordLength == 0 {
		return fmt.Errorf("provide non zero SC_MIN_WORD_LENGTH")
	}
	if o.MinWordFreq == 0 {
		return fmt.Errorf("provide non zero SC_MIN_WORD_FREQ")
	}
	if o.UnigramWeight == 0 {
		return fmt.Errorf("provide non zero SC_UNIGRAM_WEIGHT")
	}
	if o.BigramWeight == 0 {
		return fmt.Errorf("provide non zero SC_BIGRAM_WEIGHT")
	}
	if o.TrigramWeight == 0 {
		return fmt.Errorf("provide non zero SC_TRIGRAM_WEIGHT")
	}
	return nil
}

func New() (*Config, error) {
	var cfg Config
	loader := env.New("SC_", log.Printf)
	if err := loader.Load(&cfg); err != nil {
		return nil, err
	}
	if cfg.Addr == "" {
		cfg.Addr = ":10000"
	}
	if cfg.MinWordLength == 0 {
		cfg.MinWordLength = 2
	}
	if cfg.MinWordFreq == 0 {
		cfg.MinWordFreq = 5
	}
	if cfg.UnigramWeight == 0 {
		cfg.UnigramWeight = 100
	}
	if cfg.BigramWeight == 0 {
		cfg.BigramWeight = 15
	}
	if cfg.TrigramWeight == 0 {
		cfg.TrigramWeight = 5
	}
	return &cfg, cfg.Validate()
}
