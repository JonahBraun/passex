package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	
	// web
	"net/http"
	"html/template"
)

var (
	verbose = flag.Bool("v", false, "verbose") 
	
	alphabet = map[string] []rune {
		"easy": []rune("abcdefghjkmnpqrstuvwxqz23456789"),
		"full": []rune("abcdefghijklmnopqrstuvwxqzABCDEFGHIJKLMNOPQRSTUVWXQZ0123456789"),
		"punc": []rune("abcdefghijklmnopqrstuvwxqzABCDEFGHIJKLMNOPQRSTUVWXQZ0123456789~!@#$%^&*()-_+={}[]|;:<>?,./"),
		"pepr": []rune("ABCDEFGHIJKLMNOPQRSTUVWXQZ0123456789~!@#$%^&*()-_+={}[]|;:<>?,./"),
	}

	runeChan = make(chan rune, 2048)
)


func handler(w http.ResponseWriter, r *http.Request) {
	h := r.URL.Path[1:5]
	fmt.Fprintf(w, "Hi there, I love %s!", h)
	t, _ := template.ParseFiles("frontpage.html")
	t.Execute(w, nil)
}

func genRune(c chan rune, set []rune) {
	for {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		
		c <- set[random.Int64()]
	}
}

func makePass(w http.ResponseWriter, req *http.Request){
	log.Print(req.URL)
	return
	
	length := 4
	runes := make([]rune, length)

	for i := range runes { runes[i] = <- runeChan }

	pass := string(runes)
	fmt.Println(pass)
	io.WriteString(w, pass)
}

func main() {
	flag.Parse()

	log.Print("Starting up...")

	go genRune(runeChan, alphabet["punc"])

	http.HandleFunc("/pass/", makePass)
	http.ListenAndServe(":8080", nil)
}
