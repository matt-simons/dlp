package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

func main() {
	//go memleak()

	port := ":8000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.SetOutput(os.Stdout)

	logLevel, _ := log.ParseLevel(os.Getenv("LOG"))
	log.SetLevel(logLevel)

	log.WithFields(log.Fields{
		"LOG": logLevel,
	}).Info("LogLevel set")

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
