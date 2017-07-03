package main

import (
	"fmt"
	"google.golang.org/appengine"
	"http"
	//	_ "myapp/package0"
	//	_ "myapp/package1"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello world!")
}

func main() {

	http.HandleFunc("/", handler)

	appengine.Main()
}
