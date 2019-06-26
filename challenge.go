package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	content, Rerr := ioutil.ReadFile("data.txt")
	if Rerr != nil {
		log.Fatal(Rerr)
	}

	fmt.Printf("File contents: %s\n", content)

	i, err := strconv.ParseInt(strings.TrimSpace(string(content[:])), 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("Incremented file contents: %s\n", strconv.FormatInt(i+1, 16))

	message := []byte(strconv.FormatInt(i+1, 10))
	//message := []byte("4")

	Werr := ioutil.WriteFile("data.txt", message, 0644)

	if Werr != nil {
		log.Fatal(Werr)
	}
	//*/
}
