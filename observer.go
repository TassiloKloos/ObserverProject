package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

var availableApps []string
var availableAppsButtonHTML string
var runningProcesses []*exec.Cmd

func readXML() { //später hier das XML auslesen
	availableApps = append(availableApps, "C:/Program Files (x86)/Mozilla Firefox/firefox.exe")
	availableApps = append(availableApps, "C:/WINDOWS/system32/notepad.exe")
	a := "<form action='/procStart/?procStartID=0' method='post'><input type='submit' value='Firefox Starten'></form>"
	b := "<form action='/procStart/?procStartID=1' method='post'><input type='submit' value='Notepad Starten'></form>"
	availableAppsButtonHTML = a + b
}

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

func procStartHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	procStart, _ := strconv.ParseInt(q.Get("procStartID"), 0, 32)

	cmd := exec.Command(availableApps[procStart])
	cmd.Run()
	runningProcesses = append(runningProcesses, cmd)

	//evtl hier Fehler abfangen
	//stdout, err := cmd.Output()
	//if err != nil {
	//	println(err.Error())
	//}
	//print(string(stdout))

	responseString := "<html><head><title></title></head><body>" +
		"<form action='/proc/?procNr=ID0' method='post'><input type='submit' value='Prozess 0'></form>" +
		"<form action='/proc/?procNr=ID1' method='post'><input type='submit' value='Prozess 1'></form>" +
		availableAppsButtonHTML +
		"</body></html>"
	w.Write([]byte(responseString))
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
		availableAppsButtonHTML +
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
	readXML()

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/proc/", procHandler)
	http.HandleFunc("/procStart/", procStartHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))

}
