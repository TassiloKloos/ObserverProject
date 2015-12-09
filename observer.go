package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
)

var availableApps = make([]string, 0)
var stdinPipes = make([]io.WriteCloser, 0)
var availableAppsButtonHTML string
var runningProcesses = make([]*exec.Cmd, 0)
var stopProcessButtons string
var restartCounter = make([]int, 0)
var automaticRestart = make([]bool, 0)
var operatingSystem string

const responseStringFirstLine string = "<html><head><title></title></head><body>"
const responseStringLastLine string = "</body></html>"

func readXML() { //später hier das XML auslesen
	availableApps = append(availableApps, "C:/Program Files (x86)/Mozilla Firefox/firefox.exe")
	availableApps = append(availableApps, "C:/WINDOWS/system32/notepad.exe")
	availableApps = append(availableApps, "D:\\Uni\\5. Semester\\Programmieren 2\\test.bat")

	for procStartID := range availableApps { //dynamische Erzeugung der Buttons
		availableAppsButtonHTML = availableAppsButtonHTML + "<form action='/procStart/?procStartID=" + strconv.Itoa(procStartID) + "&autoRestart=false' method='post'><input type='submit' value='" + availableApps[procStartID] + " Starten'></form>" +
			"<form action='/procStart/?procStartID=" + strconv.Itoa(procStartID) + "&autoRestart=true' method='post'><input type='submit' value='" + availableApps[procStartID] + " Starten (mit automatischem Neustart)'></form>"
	}

	for i := 0; i < len(availableApps); i++ {
		restartCounter = append(restartCounter, 0)
		//		automaticRestart = append(automaticRestart, false)
	}
}

func createStopButtons() {
	stopProcessButtons = ""
	for processNR := range runningProcesses { //überprüfen, ob Prozess noch läuft; automaticRestart/restartCounter beachten
		//		_, err := os.FindProcess(runningProcesses[processNR].Process.Pid)
		//fmt.Println(proc.Pid)

		//		channel := make(chan error, 1)
		//		go func() {
		//			channel <- runningProcesses[processNR].Wait()
		//		}()
		//		err := <-channel
		//		if err != nil {
		//			//		if runningProcesses[processNR].ProcessSta != nil { //still not working
		stopProcessButtons = stopProcessButtons + "<form action='/procKill/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " killen'></form>"
		//		}
	}
}

func procStartHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	procStart, _ := strconv.ParseInt(q.Get("procStartID"), 0, 32)

	cmd := exec.Command(availableApps[procStart])

	//cmd := exec.Command("D:\\Uni\\5. Semester\\Programmieren 2\\test.bat")

	if stdin, err := cmd.StdinPipe(); err != nil {
		log.Fatal(err)
		stdinPipes = append(stdinPipes, stdin)
	}
	cmd.Start()
	if q.Get("autoRestart") == "true" {
		automaticRestart = append(automaticRestart, true)
	} else {
		automaticRestart = append(automaticRestart, false)
	}
	//fmt.Println(cmd.Process.Pid) //noch zu entfernen
	runningProcesses = append(runningProcesses, cmd)

	//	cmd.Process.Kill()
	createStopButtons()
	responseString := responseStringFirstLine +
		availableAppsButtonHTML +
		"<form action='/proc/?procNr=ID0' method='post'><input type='submit' value='Prozess 0'></form>" +
		"<form action='/proc/?procNr=ID1' method='post'><input type='submit' value='Prozess 1'></form>" +
		stopProcessButtons +
		responseStringLastLine

	w.Write([]byte(responseString))
}

func procKillHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	procNr := q.Get("procNr")
	if procNr == "" {
		procNr = "notFound"
	}
	procToKill, err := strconv.Atoi(procNr)
	if err != nil {
		fmt.Println(err)
	} else {
		runningProcesses[procToKill].Process.Kill() //weiches beenden folgt noch
		//		cmd := exec.Command("Taskkill", "/PID", strconv.Itoa(runningProcesses[procToKill].Process.Pid), "/F")
		//		cmd.Run()
		//		defer stdinPipes[procToKill].Close()
		//		io.WriteString(stdinPipes[procToKill], "4")
	}

	createStopButtons()
	responseString := responseStringFirstLine +
		availableAppsButtonHTML +
		"<form action='/proc/?procNr=ID0' method='post'><input type='submit' value='Prozess 0'></form>" +
		"<form action='/proc/?procNr=ID1' method='post'><input type='submit' value='Prozess 1'></form>" +
		stopProcessButtons +
		responseStringLastLine
	w.Write([]byte(responseString))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	name := q.Get("name")
	if name == "" {
		name = "World"
	}
	createStopButtons()
	responseString := responseStringFirstLine +
		availableAppsButtonHTML +
		"<form action='/proc/?procNr=ID0' method='post'><input type='submit' value='Prozess 0'></form>" +
		"<form action='/proc/?procNr=ID1' method='post'><input type='submit' value='Prozess 1'></form>" +
		stopProcessButtons +
		responseStringLastLine
	w.Write([]byte(responseString))
}

func main() {
	readXML()

	if runtime.GOOS == "windows" {
		operatingSystem = "windows"
	}

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/procKill/", procKillHandler)
	http.HandleFunc("/procStart/", procStartHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))

}
