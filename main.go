package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	//go memleak()

	// Use aternative port if env variable is set
	port := ":8000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// Start server
	router := NewRouter()
	log.Fatal(http.ListenAndServe(port, router))
}

//func memleak() {
//	var s [][]int
//	for {
//		s = append(s, make([]int, 100))
//		fmt.Printf("leaking memory %s", s[0])
//		time.Sleep(100 * time.Millisecond)
//	}
//}
