package main

import (
	"flag"
	//"fmt"
	"io"
	"log"
	"strconv"

	// web
	"html/template"
	"github.com/bmizerany/pat"
	"net/http"
)

type PassConfig struct{
	alphabet string
	Pass func(length int) ([]rune)
}

var (
	verbose = flag.Bool("v", false, "verbose")

	defaultFounts = map[string] PassConfig{
		"lowercase": {"abcdefghjkmnpqrstuvwxqz", func(length int, f *Fount)(r []rune){
			return []rune("foo")
		}},
		//"easy": "abcdefghjkmnpqrstuvwxqz23456789",
		//"full": "abcdefghijklmnopqrstuvwxqzABCDEFGHIJKLMNOPQRSTUVWXQZ0123456789",
		//"puncuation": "abcdefghijklmnopqrstuvwxqzABCDEFGHIJKLMNOPQRSTUVWXQZ0123456789~!@#$%^&*()-_+={}[]|;:<>?,./",
		//"pepper": "ABCDEFGHIJKLMNOPQRSTUVWXQZ0123456789~!@#$%^&*()-_+={}[]|;:<>?,./",
		//"consonants": "bcdfghjklmnpqrstvwxz",
		//"vowels": "aeiouy",
	}

	founts = map[string] *Fount{}
)

func page(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("frontpage.html")
	t.Execute(w, nil)
}

/*
func fakeWord(c chan string){
	for {
		length, err := rand.Int(rand.Reader, big.NewInt(int64(5)))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// the length of words will be between 2 and 7 characters
		word := make([]rune, length+2)
		
		state := ""
		for i := range word{
			switch state{
			case "":
				if rand.Int(rand.reader, big.NewInt(1)) {
					state = "c"
					word[i] = <- rc["consonants"]
				}
				else {
					state = "v"
					word[i] = <- rc["vowels"]
				}
*/

func makePass(w http.ResponseWriter, req *http.Request) {
	length, err := strconv.Atoi(req.URL.Query().Get(":length"))

	if err != nil || length>2000 {
		http.Error(w, "Length must be an int < 2000", http.StatusBadRequest);
		return
	}

	fountType := req.URL.Query().Get(":type")
	fount, ok := founts[fountType];

	if !ok {
		http.Error(w, "Not a valid pass type", http.StatusBadRequest)
		return
	}

	runes := fount.runes(length)

	pass := string(runes)
	log.Print("Generated password: " + pass)
	io.WriteString(w, pass)
}

func main() {
	flag.Parse()
	log.Print("PassEx web server starting")

	log.Print("Creating rune fountains...")
	for i, config := range defaultFounts {
		founts[i] = NewFount(config)
	}

	m := pat.New()
	m.Get("/make/:type/:length", http.HandlerFunc(makePass))

	m.Get("/", http.HandlerFunc(page))

	http.Handle("/", m)
	http.ListenAndServe(":8080", nil)
}

