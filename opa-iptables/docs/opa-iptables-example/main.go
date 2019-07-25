package main

import (
	"log"
	"fmt"
	"net/http"
	"os"
)

var (
	port = ""
)

func getEnv(env string) string {
	e := os.Getenv(env)
	return e
}

func main() {
	port = getEnv("PORT")
	if port == "" {
		port = "9090"
	}
	fmt.Println("Server is running on port:",port)
	http.HandleFunc("/",rootHandler)
	err := http.ListenAndServe(":"+port,nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"%s\n%s","Server is running on port: "+port,"Hello world!!")
}