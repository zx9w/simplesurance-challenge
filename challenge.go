package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
	// TODO: Stress tests and accuracy tests
)

func TestServer(*testing.T) {

}

func BenchmarkServer(*testing.B) {

}

func main() {
	// These functions are the entrypoint
	http.HandleFunc("/", SolutionServer)
	http.ListenAndServe(":8080", nil)

	// TODO: Read how to have them update a datastructure I can listen to here.

	// TODO: Async thread that has access to same datastructure (read only).

	// TODO: Efficient writing to file when datastructure meets predicate.
	// ----> This means figuring out the correct predicate
	// ----> Best datastructure (see issue on time layout)
	// ----> ----> Indexed Queue (scheduler)
	// ----> Parsing and Writing
}

func ReadData(filename string) []byte {
	content, Rerr := ioutil.ReadFile(filename)
	if Rerr != nil {
		log.Fatal(Rerr)
	}
	return content
}

func WriteData(filename string, message []byte) {
	Werr := ioutil.WriteFile(filename, message, 0644)
	if Werr != nil {
		log.Fatal(Werr)
	}
}

func SolutionServer(w http.ResponseWriter, r *http.Request) {
	// TODO: Study how I can make layout more efficient, better precision.
	// ----> This means solving for what information needs to be kept.
	// ----> Finding the notation to express that intent.
	layout := "Mon Jan 2 15:04:05 MST 2006  (MST is GMT-0700)"

	content := ReadData("date.txt")
	//	fmt.Printf("Dates:\n%s\n", content)

	lines := strings.Split(string(content[:]), "\n")

	nowt := time.Now()
	//	fmt.Printf("Time now:\n-> %s\n\n", nowt.Format(layout))

	counter := 0
	write := nowt.Format(layout)
	for line := 0; line < len(lines); line++ {
		oldt, _ := time.Parse(layout, lines[line])
		cmpt := oldt.Add(1 * time.Minute)

		if nowt.Before(cmpt) {
			write = write + "\n" + oldt.Format(layout)
			counter += 1
		} else {
			//			fmt.Printf("Discarded: %s\n", oldt.Format(layout))
		}
	}

	fmt.Fprintf(w, "Number of calls in last 60 seconds: %s", strconv.Itoa(counter))

	WriteData("date.txt", []byte(write))

}
