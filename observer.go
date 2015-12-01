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
	responseString := "<html><head><title></title></head><body><form action='/login' method='post'><input type='submit' value='Print asdf'></form></body></html>"
	w.Write([]byte(responseString))
}
func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("asdf")
	responseString := "<html><head></head><body>asdf</body></html>"
	w.Write([]byte(responseString))
}
func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/asdf", mainHandler)
	http.HandleFunc("/login", testHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
