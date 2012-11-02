/*
	Package fount provides:
	- Fount: random rune fountains of the provided alphabet.
	- Pass: Password generator
	
*/
package main

import (
	"crypto/rand"
	"log"
	"math/big"
)

type Fount struct{
	Name string
	alphabet []rune
	fount chan rune
	Pass func(length int) ([]rune)
}

func NewFount(alphabet string) (f *Fount) {
	f = &Fount{
		alphabet: []rune(alphabet),
		fount: make(chan rune),
	}
	go f.run()
	return
}

func (f *Fount) run () {
	for {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(len(f.alphabet))))
		if err != nil {
			log.Fatal(err)
		}

		f.fount <- f.alphabet[random.Int64()]
	}
}

func (f *Fount) runes(length int) (r []rune) {
	r = make([]rune, length)
	
	for i := range r {
		r[i] = <-f.fount
	}
	return
}
