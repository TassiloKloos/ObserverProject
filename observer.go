package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func openShell(list string) bool {
	cmd := exec.Command(list)
	stdout, err := cmd.Output()
	if err != nil {
		//fmt.Println(err)
		return false
	}
	fmt.Println(string(stdout))
	return true
}

func enterCommand(command string, location string) bool {
	cmd := exec.Command(command, location)

	cmd.Run()
	stdout, err := cmd.Output()
	if !strings.Contains(err.Error(), "asdf") {
		//fmt.Println(err)
		return false
	}
	fmt.Println(string(stdout))
	return true
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	name := q.Get("name")
	if name == "" {
		name = "World"
	}
	responseString := "<html><head><title></title></head><body>" +
		"<form action='/proc/?procNr=ID0' method='post'><input type='submit' value='Prozess 0'></form>" +
		"<form action='/proc/?procNr=ID1' method='post'><input type='submit' value='Prozess 1'></form>" +
		"</body></html>"
	w.Write([]byte(responseString))
}
func procHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	procNr := q.Get("procNr")
	if procNr == "" {
		procNr = "notFound"
	}
	fmt.Println(procNr)
	responseString := "<html><head></head><body>" + procNr + "</body></html>"
	w.Write([]byte(responseString))
}
func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/proc/", procHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
