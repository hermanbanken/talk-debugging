package main

import (
	"net/http"
	"os"
)

func main() {
	http.ListenAndServe(addr(), http.DefaultServeMux)

}

func addr() string {
	if port, hasPort := os.LookupEnv("PORT"); hasPort {
		return ":" + port
	}
	return ":8080"
}
