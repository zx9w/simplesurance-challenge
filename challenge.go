package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	// ----- Read from file -----
	content, Rerr := ioutil.ReadFile("data.txt")
	if Rerr != nil {
		log.Fatal(Rerr)
	}

	fmt.Printf("File contents: %s\n", content)

	// ------ Type conversions ----
	i, Cerr := strconv.ParseInt(strings.TrimSpace(string(content[:])), 10, 64)
	if Cerr != nil {
		log.Fatal(Cerr)
	}

	message := []byte(strconv.FormatInt(i+1, 10))

	// ------ Write to file ------
	Werr := ioutil.WriteFile("data.txt", message, 0644)
	if Werr != nil {
		log.Fatal(Werr)
	}

	// ------ HTTP ------
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)

	//*/
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
