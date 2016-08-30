package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 6)
	fmt.Println("gomaxprocs is set to", runtime.GOMAXPROCS(0))
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/playerstat", func(w http.ResponseWriter, r *http.Request) {
		// Print player's browser information to stdout
		log.Printf("Player %s (%s x %s) is a '%s'", r.RemoteAddr, r.FormValue("width"), r.FormValue("height"), r.Header.Get("User-Agent"))
	})
	http.HandleFunc("/playbackstat", func(w http.ResponseWriter, r *http.Request) {
		// Print playback information submitted by browser to stdout
		log.Printf("Player %s is playing %s (%s x %s) at %s", r.RemoteAddr, r.FormValue("currenttype"), r.FormValue("width"), r.FormValue("height"), r.FormValue("currenttime"))
	})
	http.HandleFunc("/randvideo", func(w http.ResponseWriter, r *http.Request) {
		// Serve a random video among 99 video files. cpsamples.sh makes the 99 copies.
		http.ServeFile(w, r, "sample-video-"+strconv.Itoa(1+rand.Intn(98))+".webm")
	})
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Serve HTTPS on 10987 and HTTP on 10988
	go func() {
		if err := http.ListenAndServeTLS("0.0.0.0:10987", "tls.crt", "tls.key", nil); err != nil {
			panic(err)
		}
	}()
	if err := http.ListenAndServe("0.0.0.0:10988", nil); err != nil {
		panic(err)
	}
}
