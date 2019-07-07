package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	//"strings"
	//"testing"
	"time"
)

func mkByte(t time.Time, layout string) []byte { // FIXME
	return []byte(t.Format(layout))

}

func Funnel(q *list.List, reqs <-chan time.Time, l string, write1, write2 chan<- []byte) {
	que := *q
	for {
		select {
		case timestamp := <-reqs:
			// These write to files.
			write1 <- mkByte(timestamp, l)
			write2 <- mkByte(timestamp, l)

			que.PushFront(timestamp)
			for e := que.Back(); e != nil && e.Value.(time.Time).Add(1*time.Minute).Before(timestamp); e = e.Prev() {
				que.Remove(e)
				e = que.Back()
			}
		default: // nothing to write
			for e := que.Back(); e != nil && e.Value.(time.Time).Add(1*time.Minute).Before(time.Now()); e = e.Prev() {
				que.Remove(e)
			}
			// consider sleeping for some nanoseconds
		}
	}
}

func Write(file string, write <-chan []byte, stop <-chan bool, resp chan<- bool) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	running := true

	for running {
		select {
		case msg := <-write:
			w.Write(msg)
		case stop := <-stop:
			running = !stop
		default:
			w.Flush()
		}
	}
	w.Flush()
	resp <- true
}

func Clean(write1, write2 chan []byte) {

	stop2 := make(chan bool)
	stop1 := make(chan bool)
	resp2 := make(chan bool)
	resp1 := make(chan bool)

	go Write("data1.txt", write1, stop1, resp1) // true
	go Write("data2.txt", write2, stop2, resp2) // false

	overwrite := false
	for {
		time.Sleep(70 * time.Second)
		if overwrite {
			stop1 <- true
			if <-resp1 {
				WriteData("data1.txt", nil)
				go Write("data1.txt", write1, stop1, resp1)
			}
		} else {
			stop2 <- true
			if <-resp2 {
				WriteData("data2.txt", nil)
				go Write("data2.txt", write2, stop2, resp2)
			}
		}
		overwrite = !overwrite
	}
}

func Init(que *list.List, layout string) {
	data1, err1 := os.Open("data1.txt")
	data2, err2 := os.Open("data2.txt")
	if err1 != nil && err2 != nil {
		log.Fatal(err1)
		log.Fatal(err2)
	}

	defer data1.Close()
	defer data2.Close()

	scan1 := bufio.NewScanner(data1)
	scan2 := bufio.NewScanner(data2)

	var l1, l2 time.Time
	if scan1.Scan() {
		l1, _ = time.Parse(layout, scan1.Text())
	} else {
		l1 = time.Now()
	}
	if scan2.Scan() {
		l2, _ = time.Parse(layout, scan2.Text())
	} else {
		l2 = time.Now()
	}

	var scan *bufio.Scanner
	if l1.Before(l2) {
		scan = bufio.NewScanner(data1)
	} else {
		scan = bufio.NewScanner(data2)
	}

	nowt := time.Now()
	queue := *que

	for scan.Scan() {
		oldt, _ := time.Parse(layout, scan.Text())
		cmpt := oldt.Add(1 * time.Minute)

		if nowt.Before(cmpt) {
			queue.PushFront(oldt)
		}
	}
}

func main() {
	layout := "Mon Jan 2 15:04:05 MST 2006  (MST is GMT-0700)"
	queue := list.New()

	Init(queue, layout)

	reqs := make(chan time.Time, 300) // sends reqs to funnel
	write2 := make(chan []byte, 300)  // writes to data2.txt
	write1 := make(chan []byte, 300)  // data1.txt

	go Funnel(queue, reqs, layout, write1, write2)
	go Clean(write1, write2)

	// These functions are the entrypoint
	http.HandleFunc("/", Solution(reqs, queue))
	http.ListenAndServe(":8080", nil)
}

func WriteData(filename string, message []byte) {
	Werr := ioutil.WriteFile(filename, message, 0644)
	if Werr != nil {
		log.Fatal(Werr)
	}
}

func Solution(store chan<- time.Time, q *list.List) func(w http.ResponseWriter, r *http.Request) {

	log.Printf("I was here 1")

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Requests in past 60 seconds: %s", strconv.Itoa((*q).Len()+1))
		store <- time.Now()
	}
}
