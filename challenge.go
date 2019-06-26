package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	// ------ HTTP ------
	http.HandleFunc("/", SolutionServer)
	http.ListenAndServe(":8080", nil)

	//*/
}

// ----- Read from file -----> Side effects!
func ReadData(filename string) []byte {
	content, Rerr := ioutil.ReadFile(filename)
	if Rerr != nil {
		log.Fatal(Rerr)
	}
	return content
}

// ------ Write to file ------> Side effects!
func WriteData(filename string, message []byte) {
	Werr := ioutil.WriteFile(filename, message, 0644)
	if Werr != nil {
		log.Fatal(Werr)
	}
}

// ------ HTTP Server ------> Side effects!
func SolutionServer(w http.ResponseWriter, r *http.Request) {
	layout := "Mon Jan 2 15:04:05 MST 2006  (MST is GMT-0700)"

	content := ReadData("date.txt")
	fmt.Printf("Dates:\n %s\n", content)

	nowt := time.Now()
	fmt.Printf("Time now:\n %s\n", nowt.Format(layout))

	oldt, _ := time.Parse(layout, string(content[:]))
	fmt.Printf("Time from file, before increment:\n %s\n", oldt.Format(layout))
	oldt = oldt.Add(1 * time.Minute)
	fmt.Printf("Time from file, after increment:\n %s\n", oldt.Format(layout))

	counter := 0

	if nowt.Before(oldt) {
		fmt.Printf("The last call will be counted.\n")
		counter += 1
	} else {
		fmt.Printf("The last call won't be counted.\n")
	}

	fmt.Fprintf(w, "Number of calls in last 60 seconds: %s", strconv.Itoa(counter))

	WriteData("date.txt", []byte(nowt.Format(layout)))

}

// ------ HTTP Server ------> Side effects!
func HelloServer(w http.ResponseWriter, r *http.Request) {
	content := ReadData("data.txt")
	fmt.Printf("File contents: %s\n", content)

	fmt.Fprintf(w, "File contents: %s!", content) //r.URL.Path[1:])

	// ------ Type conversions ----
	i, Cerr := strconv.ParseInt(strings.TrimSpace(string(content[:])), 10, 64)
	if Cerr != nil {
		log.Fatal(Cerr)
	}

	message := []byte(strconv.FormatInt(i+1, 10))
	// ============================

	WriteData("data.txt", message)

}
