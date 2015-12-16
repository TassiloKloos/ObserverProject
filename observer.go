package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

var availableApps = make([]string, 0)
var stdinPipes = make([]io.WriteCloser, 0)
var availableAppsButtonHTML string
var runningProcesses = make([]*exec.Cmd, 0)
var stopProcessButtons string
var restartCounter = make([]int, 0)
var automaticRestart = make([]bool, 0)
var counterHTML string
var outputButtonHTML string
var ShellOutput = make([]string, 0)

const responseStringFirstLine string = "<html><head><title></title></head><body>"
const responseStringLastLine string = "</body></html>"

func readXML() { //später hier das XML auslesen
	//availableApps = append(availableApps, "C:/Program Files (x86)/Mozilla Firefox/firefox.exe")
	//availableApps = append(availableApps, "C:/WINDOWS/system32/notepad.exe")
	availableApps = append(availableApps, "D:\\Uni\\5. Semester\\Programmieren 2\\test.bat")
	availableApps = append(availableApps, "D:\\Uni\\5. Semester\\Programmieren 2\\test2.bat")

	for procStartID := range availableApps { //dynamische Erzeugung der Buttons
		availableAppsButtonHTML = availableAppsButtonHTML + "<form action='/procStart/?procStartID=" + strconv.Itoa(procStartID) + "&autoRestart=false' method='post'><input type='submit' value='" + availableApps[procStartID] + " Starten'></form>" +
			"<form action='/procStart/?procStartID=" + strconv.Itoa(procStartID) + "&autoRestart=true' method='post'><input type='submit' value='" + availableApps[procStartID] + " Starten (mit automatischem Neustart)'></form>"
	}

	for i := 0; i < len(availableApps); i++ {
		restartCounter = append(restartCounter, 0)
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
		stopProcessButtons = stopProcessButtons + "<form action='/procKill/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " hart beenden'></form>" +
			"<form action='/procSoftKill/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " weich beenden'></form>"
		//		}
	}
}

func createShellOutputButtons() {
	outputButtonHTML = ""
	for appNR := range runningProcesses {
		outputButtonHTML = outputButtonHTML + "<form action='/procShellOutput/?procNr=" + strconv.Itoa(appNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[appNR].Path + " Shell Output zeigen'></form>"
	}
}

func createCounter() {
	counterHTML = ""
	for appNR := range availableApps {
		counterHTML = counterHTML + "<label>" + availableApps[appNR] + " wurde " + strconv.Itoa(restartCounter[appNR]) + " mal neu gestartet</label><br><br>"
	}
}

func procStartHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	procStart, _ := strconv.ParseInt(q.Get("procStartID"), 0, 32)

	cmd := exec.Command(availableApps[procStart])

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	} else {
		stdinPipes = append(stdinPipes, stdin)
	}

	//	stdout, err := cmd.StdoutPipe()
	//	if err != nil {
	//		log.Fatal(err)
	//	} else {
	//		//		ShellOutput = append(ShellOutput, asdf)
	//	}

	//	writer := bufio.ReadWriter(cmd.Stdin, cmd.Stdout)
	//	writer.WriteString
	//	defer writer.Flush()
	cmd.Start()

	if q.Get("autoRestart") == "true" {
		automaticRestart = append(automaticRestart, true)
	} else {
		automaticRestart = append(automaticRestart, false)
	}

	runningProcesses = append(runningProcesses, cmd)

	createCounter()
	createStopButtons()
	createShellOutputButtons()
	responseString := responseStringFirstLine +
		availableAppsButtonHTML +
		counterHTML +
		stopProcessButtons +
		outputButtonHTML +
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
		runningProcesses[procToKill].Process.Kill()
		automaticRestart[procToKill] = false
	}
	createCounter()
	createStopButtons()
	createShellOutputButtons()
	responseString := responseStringFirstLine +
		availableAppsButtonHTML +
		counterHTML +
		stopProcessButtons +
		outputButtonHTML +
		responseStringLastLine
	w.Write([]byte(responseString))
}

func procSoftKillHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	procNr := q.Get("procNr")
	if procNr == "" {
		procNr = "notFound"
	}
	procToKill, err := strconv.Atoi(procNr)
	if err != nil {
		fmt.Println(err)
	} else {
		writer := bufio.NewWriter(stdinPipes[procToKill])
		writer.WriteString("stop")
		writer.Flush()
		automaticRestart[procToKill] = false
	}
	createCounter()
	createStopButtons()
	createShellOutputButtons()
	responseString := responseStringFirstLine +
		availableAppsButtonHTML +
		counterHTML +
		stopProcessButtons +
		outputButtonHTML +
		responseStringLastLine
	w.Write([]byte(responseString))
}

func procShellOutputHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	procNr := q.Get("procNr")
	if procNr == "" {
		procNr = "notFound"
	}
	procToOutput, err := strconv.Atoi(procNr)
	if err != nil {
		fmt.Println(err)
	}

	//abfrage ob prozess noch läuft (siehe Stop Buttons)
	//	stdout, err := ioutil.ReadAll(runningProcesses[procNr].Stdout)
	//	if err != nil {
	//		log.Fatal(err)
	//	} else {
	//		ShellOutput[procNr] = string(stdout)
	//	}

	createCounter()
	createStopButtons()
	createShellOutputButtons()
	responseString := responseStringFirstLine +
		"<script>alert(" + ShellOutput[procToOutput] + ")</script>" +
		availableAppsButtonHTML +
		counterHTML +
		stopProcessButtons +
		outputButtonHTML +
		responseStringLastLine
	w.Write([]byte(responseString))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	createCounter()
	createStopButtons()
	createShellOutputButtons()
	responseString := responseStringFirstLine +
		availableAppsButtonHTML +
		counterHTML +
		stopProcessButtons +
		outputButtonHTML +
		responseStringLastLine
	w.Write([]byte(responseString))
}

func main() {
	readXML()

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/procKill/", procKillHandler)
	http.HandleFunc("/procSoftKill/", procSoftKillHandler)
	http.HandleFunc("/procStart/", procStartHandler)
	http.HandleFunc("/procShellOutput/", procShellOutputHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))

}
