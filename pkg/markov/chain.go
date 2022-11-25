package markov

import (
	"fmt"
	"math/rand"
	"sync"
	"unicode/utf8"
)

type Chain struct {
	mu     sync.Mutex
	ran    *rand.Rand
	order  int
	ngrams map[string][]rune
}

func NewChain(order int, randsrc rand.Source) *Chain {
	if order <= 0 {
		order = 5
	}
	return &Chain{
		order:  order,
		ran:    rand.New(randsrc),
		ngrams: make(map[string][]rune),
	}
}

func (c *Chain) Learn(corpus string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Range over corpus n-grams
	corpusChars := []rune(corpus)
	for i := 0; i < len(corpusChars)-c.order; i++ {
		// Add next character to ngrams for current
		currSequence := string(corpusChars[i : i+c.order])
		c.ngrams[currSequence] = append(c.ngrams[currSequence], corpusChars[i+c.order])
	}
}

func (c *Chain) NextCharacter(str string) rune {
	strNumChars := utf8.RuneCountInString(str)
	if strNumChars < c.order {
		panic(fmt.Errorf("input string %d is shorter than order %d", strNumChars, c.order))
	}
	lastNchars := []rune(str)[strNumChars-c.order:]
	options, ok := c.ngrams[string(lastNchars)]
	if !ok {
		panic(fmt.Errorf("unknown ngram: %q", lastNchars))
	}
	return options[c.ran.Intn(len(options))]
}

func (c *Chain) String() string {
	out := ""
	for k, v := range c.ngrams {
		out += fmt.Sprintf("%q: %q\n", k, v)
	}
	return out
}
