# Spell correct

Performs context based spelling correction.

Based on: https://arxiv.org/pdf/1910.11242.pdf

work in progress. NOT yet production ready

## How to use

See the most important environment variables in order to get started.

```
SC_ADDR: The address the server listens to (default :10000)
SC_SENTENCES_PATH: The path of the file containing sentences from the language
                    you want to perform spelling correction for. (File must be gzipped)
SC_DICT_PATH: The path of the file  with word frequency dictionary for the languate. (file must be gzipped)
```

See the internal/config/config.go for additional variables

Example assuming that the train files are in the datasets/ folder.


1. `make build`
2. `SC_ADDR=:10000 SC_SENTENCES_PATH=datasets/en/sentences.txt.gz SC_DICT_PATH=datasets/en/freq-dict.txt.gz ./spell-correct-server`

It starts a web server listening by default on port 10.000


```
curl 'http://localhost:10000/?query=scred%20a%20bicycle%20kikc'
```

It gives you back the suggestions. The suggestions are ordered with the ones the algorithm decides are most 
relevant first.

### Using docker

Example:

```
make docker-build
docker run  -p 10000:10000 -v /home/gosom/datasets:/datasets -e SC_SENTENCES_PATH=/datasets/en/sentences.txt.gz -e SC_DICT_PATH=/datasets/en/freq-dict.txt.gz spell-correctort
```


## Supported Languages

The method is language independent.
All you need is two files used for training:

- file containing sentences for the languages
- file containing the word frequencies for the languages

In the datasets directory we added examples for English and German.

Datasets source is referenced in the README.md in the language's folder


#### Special Thanks

- https://github.com/Oikopedo

- https://arxiv.org/pdf/1910.11242.pdf

- https://github.com/wolfgarbe/SymSpell
- https://github.com/eskriett/spell
- https://github.com/json-iterator/go
- https://github.com/qiangxue/go-env

