package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

var availableApps []string
var availableAppsButtonHTML string
var runningProcesses []*exec.Cmd
var stopProcessButtons string

const responseStringFirstLine string = "<html><head><title></title></head><body>"
const responseStringLastLine string = "</body></html>"

func readXML() { //später hier das XML auslesen
	availableApps = append(availableApps, "C:/Program Files (x86)/Mozilla Firefox/firefox.exe")
	availableApps = append(availableApps, "C:/WINDOWS/system32/notepad.exe")
	availableApps = append(availableApps, "D:\\Uni\\5. Semester\\Programmieren 2\\test.bat")
	a := "<form action='/procStart/?procStartID=0' method='post'><input type='submit' value='Firefox Starten'></form>"
	b := "<form action='/procStart/?procStartID=1' method='post'><input type='submit' value='Notepad Starten'></form>"
	c := "<form action='/procStart/?procStartID=2' method='post'><input type='submit' value='TestApp Starten'></form>"
	availableAppsButtonHTML = a + b + c
}

func createStopButtons() {
	stopProcessButtons = ""
	for processNR := range runningProcesses { //überprüfen, ob Prozess noch läuft
		proc, err := os.FindProcess(runningProcesses[processNR].Process.Pid)
		fmt.Println(proc.Pid)
		if err == nil {
			stopProcessButtons = stopProcessButtons + "<form action='/proc/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " anhalten'></form>"
		}
	}
}

//func openShell(list string) bool {
//	cmd := exec.Command(list)
//	stdout, err := cmd.Output()
//	if err != nil {
//		//fmt.Println(err)
//		return false
//	}
//	fmt.Println(string(stdout))
//	return true
//}

//func enterCommand(command string, location string) bool {
//	cmd := exec.Command(command, location)

//	cmd.Run()
//	stdout, err := cmd.Output()
//	if !strings.Contains(err.Error(), "started") {
//		//fmt.Println(err)
//		return false
//	}
//	fmt.Println(string(stdout))
//	return true
//}

func procStartHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	procStart, _ := strconv.ParseInt(q.Get("procStartID"), 0, 32)

	cmd := exec.Command(availableApps[procStart])

	//cmd := exec.Command("D:\\Uni\\5. Semester\\Programmieren 2\\test.bat")

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cmd.Process.Pid)
	runningProcesses = append(runningProcesses, cmd)

	cmd.Process.Kill()
	createStopButtons()
	responseString := responseStringFirstLine +
		availableAppsButtonHTML +
		"<form action='/proc/?procNr=ID0' method='post'><input type='submit' value='Prozess 0'></form>" +
		"<form action='/proc/?procNr=ID1' method='post'><input type='submit' value='Prozess 1'></form>" +
		stopProcessButtons +
		responseStringLastLine

	w.Write([]byte(responseString))
}

func procHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	procNr := q.Get("procNr")
	if procNr == "" {
		procNr = "notFound"
	}
	fmt.Println(procNr)
	responseString := responseStringFirstLine +
		availableAppsButtonHTML +
		"<form action='/proc/?procNr=ID0' method='post'><input type='submit' value='Prozess 0'></form>" +
		"<form action='/proc/?procNr=ID1' method='post'><input type='submit' value='Prozess 1'></form>" +
		responseStringLastLine
	w.Write([]byte(responseString))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	name := q.Get("name")
	if name == "" {
		name = "World"
	}
	responseString := responseStringFirstLine +
		availableAppsButtonHTML +
		"<form action='/proc/?procNr=ID0' method='post'><input type='submit' value='Prozess 0'></form>" +
		"<form action='/proc/?procNr=ID1' method='post'><input type='submit' value='Prozess 1'></form>" +
		responseStringLastLine
	w.Write([]byte(responseString))
}

func main() {
	readXML()

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/proc/", procHandler)
	http.HandleFunc("/procStart/", procStartHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))

}
