package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	//"strings"
	//"testing"
	"time"
)

const layout = "Mon Jan 2 15:04:05 MST 2006  (MST is GMT-0700)"

var queue *list.List

func wait(rtn chan<- bool, secs int) {
	for {
		time.Sleep(time.Duration(secs) * time.Second)
		rtn <- true
	}
}

func Funnel(reqs <-chan time.Time, write1, write2 chan<- []byte) {
	okay := make(chan bool)

	go wait(okay, 1)

	for {
		select {
		case timestamp := <-reqs:
			queue.PushFront(timestamp)

			write1 <- []byte(timestamp.Format(layout) + "\n")
			write2 <- []byte(timestamp.Format(layout) + "\n")

		case <-okay:
			for e := queue.Back(); e != nil; e = e.Prev() {
				oldt, safe := e.Value.(time.Time)
				if safe && oldt.Add(1*time.Minute).Before(time.Now()) {
					log.Printf("Removed# " + oldt.Format(layout))
					queue.Remove(e)
				}
			}
		}
	}
}

func Write(file string, write <-chan []byte, stop <-chan bool, resp chan<- bool) {

	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(f)

	running := true

	for running {
		select {
		case msg := <-write:
			w.Write(msg)
		case stop := <-stop:
			running = !stop
		}
	}
	// I was getting paranoid about how to save properly.
	w.Flush()
	f.Close()
	resp <- true
}

func stopw(stop chan<- bool, resp <-chan bool) bool {
	stop <- true
	return <-resp
}

func Clean(kill chan bool, write1, write2 chan []byte) {
	clean := make(chan bool)
	okay := make(chan bool)

	stop2 := make(chan bool)
	stop1 := make(chan bool)

	resp2 := make(chan bool)
	resp1 := make(chan bool)

	go wait(clean, 65)
	go wait(okay, 5)

	go Write("data1.txt", write1, stop1, resp1) // true
	go Write("data2.txt", write2, stop2, resp2) // false

	alive := true
	overwrite := false
	save := false

	for alive {
		select {
		case <-clean:
			if overwrite {
				if stopw(stop1, resp1) {
					log.Printf("Clearing data1.txt")
					WriteData("data1.txt", nil)
					go Write("data1.txt", write1, stop1, resp1)
				}
				overwrite = !overwrite
			} else {
				if stopw(stop2, resp2) {
					log.Printf("Clearing data2.txt")
					WriteData("data2.txt", nil)
					go Write("data2.txt", write2, stop2, resp2)
				}
			}
			overwrite = !overwrite
		case <-okay:
			if save {
				if stopw(stop1, resp1) {
					log.Printf("Restarted writer 1")
					go Write("data1.txt", write1, stop1, resp1)
				}
				save = !save
			} else {
				if stopw(stop2, resp2) {
					log.Printf("Restarted writer 2")
					go Write("data2.txt", write2, stop2, resp2)
				}
				save = !save
			}
		case <-kill:
			if stopw(stop1, resp1) && stopw(stop2, resp2) {
				log.Printf("Success~!")
				alive = false
			}
		}

	}
	kill <- true
}

func Init() {
	data1, err1 := os.OpenFile("data1.txt", os.O_RDWR|os.O_CREATE, 0755)
	data2, err2 := os.OpenFile("data2.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err1 != nil && err2 != nil {
		log.Fatal(err1)
		log.Fatal(err2)
	}

	defer data1.Close()
	defer data2.Close()

	scan1 := bufio.NewScanner(data1)
	scan2 := bufio.NewScanner(data2)

	// Refactor this logic so that it performs a merge
	// If same string in both then only write it once.

	var old1, old2 time.Time
	nowt := time.Now()
	queue1 := list.New()
	queue2 := list.New()

	for scan1.Scan() {
		old1, _ = time.Parse(layout, scan1.Text())
		if old1.Add(1 * time.Minute).After(nowt) {
			queue1.PushFront(old1)
		}
	}
	for scan2.Scan() {
		old2, _ = time.Parse(layout, scan2.Text())
		if old2.Add(1 * time.Minute).After(nowt) {
			queue2.PushFront(old2)
		}
	}

	if queue1.Len() > queue2.Len() {
		queue = queue1
	} else {
		queue = queue2
	}
	log.Printf(strconv.Itoa(queue.Len()))
}

func main() {
	Init()

	reqs := make(chan time.Time, 300) // sends reqs to funnel
	write2 := make(chan []byte, 300)  // writes to data2.txt
	write1 := make(chan []byte, 300)  // data1.txt
	kill := make(chan bool)

	go Funnel(reqs, write1, write2)
	go Clean(kill, write1, write2)

	// These functions are the entrypoint
	http.HandleFunc("/", Solution(reqs))
	http.ListenAndServe(":8082", nil)

	// Don't exit until Ctrl+C
	cmd := make(chan os.Signal, 1)
	signal.Notify(cmd, os.Interrupt, syscall.SIGINT)

	log.Printf("Received %s", <-cmd)
	log.Printf("Exiting gracefully.")
	kill <- true
	<-kill
}

func WriteData(filename string, message []byte) {
	Werr := ioutil.WriteFile(filename, message, 0644)
	if Werr != nil {
		log.Fatal(Werr)
	}
}

func Solution(store chan<- time.Time) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Requests in past 60 seconds: %s", strconv.Itoa(queue.Len()+1))
		store <- time.Now()
	}
}
