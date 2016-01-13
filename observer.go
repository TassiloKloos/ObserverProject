package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var availableApps = make([]string, 0)
var stdins = make([]string, 0)
var stopCommands = make([]string, 0)
var availableAppsNumberForCmd = make([]int64, 0)
var stdinPipes = make([]io.WriteCloser, 0)
var stdOutput = make([]string, 0)
var availableAppsButtonHTML string
var runningProcesses = make([]*exec.Cmd, 0)
var stopProcessButtons string
var restartCounter = make([]int, 0)
var automaticRestart = make([]bool, 0)
var counterHTML string
var outputButtonHTML string
var path string

const responseStringFirstLine string = "<html><head><title></title></head><body>" + "<form action='/' method='post'><input type='submit' value='Seite refreshen'></form>"
const responseStringLastLine string = "</body></html>"

type Application struct {
	XMLName xml.Name `xml:"application"`
	Path    string   `xml:"path"`
	Stdin   string   `xml:"stdin"`
	Stopexe string   `xml:"stopexe"`
}

type Observer struct {
	XMLName      xml.Name      `xml:"observer"`
	Applications []Application `xml:"application"`
}

func (a Application) String() string {
	return fmt.Sprintf("path : %s - stdin : %s - stopexe : %s \n", a.Path, a.Stdin, a.Stopexe)
}

func readXML(firstStart bool, xmlPath string) {
	if firstStart { //wird nur bei Programmstart ausgeführt
		path = xmlPath

		//Testpfad D:\Uni\5. Semester\Programmieren 2\observer.xml

		xmlFile, err := os.Open(path)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer xmlFile.Close()

		XMLdata, _ := ioutil.ReadAll(xmlFile)
		var o Observer
		xml.Unmarshal(XMLdata, &o)

		availableApps = make([]string, 0)
		stdins = make([]string, 0)
		stopCommands = make([]string, 0)
		for count := range o.Applications {

			availableApps = append(availableApps, o.Applications[count].Path)
			stdins = append(stdins, o.Applications[count].Stdin)
			stopCommands = append(stopCommands, o.Applications[count].Stopexe)
		}

		availableAppsButtonHTML = ""
		for procStartID := range availableApps { //dynamische Erzeugung der Buttons
			availableAppsButtonHTML = availableAppsButtonHTML + "<form action='/procStart/?procStartID=" + strconv.Itoa(procStartID) + "&autoRestart=false' method='post'><input type='submit' value='" + availableApps[procStartID] + " Starten'></form>" +
				"<form action='/procStart/?procStartID=" + strconv.Itoa(procStartID) + "&autoRestart=true' method='post'><input type='submit' value='" + availableApps[procStartID] + " Starten (mit automatischem Neustart)'></form>"
		}

		for i := 0; i < len(availableApps); i++ {
			restartCounter = append(restartCounter, 0)
		}
	} else { // aktualisiert stdins immer, wenn stop buttons erzeugt werden
		xmlFile, err := os.Open(path)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer xmlFile.Close()

		XMLdata, _ := ioutil.ReadAll(xmlFile)
		var o Observer
		xml.Unmarshal(XMLdata, &o)

		stdins = make([]string, 0)
		for count := range o.Applications {
			stdins = append(stdins, o.Applications[count].Stdin)
		}
	}
}

func createStopButtons() {
	readXML(false, "")

	stopProcessButtons = ""
	for processNR := range runningProcesses { //überprüfen, ob Prozess noch läuft; automaticRestart/restartCounter beachten

		channel := make(chan error, 1)
		go func() {
			channel <- runningProcesses[processNR].Wait()
		}()
		select { // die Idee einen select mit Timeout auf den Channel zu machen stammt von Oliver Raum, der Code wurde ohne seine Hilfe erstellt
		case err := <-channel:
			if err != nil {
				err = nil
			}
		case <-time.After(100000000): // 0,1 sec wait
			//			wartet nicht auf beenden vom channel
		}

		//try catch Funktion hier für
		func() {
			defer func() {
				if r := recover(); r != nil {
					//Catch hier
					stopProcessButtons = stopProcessButtons + "<form action='/procKill/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " hart beenden, automatischer Neustart = " + strconv.FormatBool(automaticRestart[processNR]) + "'></form>" +
						"<form action='/procSoftKill/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " weich beenden (stdin), automatischer Neustart = " + strconv.FormatBool(automaticRestart[processNR]) + "'></form>" +
						"<form action='/procAppKill/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " durch eine Anwendung beenden, automatischer Neustart = " + strconv.FormatBool(automaticRestart[processNR]) + "'></form>"
				}
			}()
			//Try hier
			if runningProcesses[processNR].ProcessState.Exited() == true && automaticRestart[processNR] == true {
				cmd := exec.Command(availableApps[availableAppsNumberForCmd[processNR]])

				stdin, err := cmd.StdinPipe()
				if err != nil {
					log.Fatal(err)
				} else {
					stdinPipes[processNR] = stdin
				}

				stdout, err := cmd.StdoutPipe()
				if err != nil {
					log.Fatal(err)
				}
				reader := bufio.NewReader(stdout)
				cmd.Start()
				line, _, _ := reader.ReadLine()
				stdOutput = append(stdOutput, string(line))
				runningProcesses[processNR] = cmd
				restartCounter[availableAppsNumberForCmd[processNR]]++
				stopProcessButtons = stopProcessButtons + "<form action='/procKill/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " hart beenden, automatischer Neustart = " + strconv.FormatBool(automaticRestart[processNR]) + "'></form>" +
					"<form action='/procSoftKill/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " weich beenden (stdin), automatischer Neustart = " + strconv.FormatBool(automaticRestart[processNR]) + "'></form>" +
					"<form action='/procAppKill/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " durch eine Anwendung beenden, automatischer Neustart = " + strconv.FormatBool(automaticRestart[processNR]) + "'></form>"
			} else if runningProcesses[processNR].ProcessState.Exited() == false {
				stopProcessButtons = stopProcessButtons + "<form action='/procKill/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " hart beenden, automatischer Neustart = " + strconv.FormatBool(automaticRestart[processNR]) + "'></form>" +
					"<form action='/procSoftKill/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " weich beenden (stdin), automatischer Neustart = " + strconv.FormatBool(automaticRestart[processNR]) + "'></form>" +
					"<form action='/procAppKill/?procNr=" + strconv.Itoa(processNR) + "' method='post'><input type='submit' value='Prozess " + runningProcesses[processNR].Path + " durch eine Anwendung beenden, automatischer Neustart = " + strconv.FormatBool(automaticRestart[processNR]) + "'></form>"
			}
		}()
	}
}

func createShellOutputButtons() {
	var sliceOutputButtonHTML = make([]string, 0)
	outputButtonHTML = ""
	for appNR := range runningProcesses {
		sliceOutputButtonHTML = append(sliceOutputButtonHTML, "<form action='/procShellOutput/?procNr="+strconv.Itoa(appNR)+"' method='post'><input type='submit' value='Prozess "+runningProcesses[appNR].Path+" Shell Output zeigen'></form>")
	}
	for position := range sliceOutputButtonHTML {
		if position < 10 { //limit der Output Buttons auf maximal 10 Stück
			outputButtonHTML = outputButtonHTML + sliceOutputButtonHTML[len(sliceOutputButtonHTML)-1-position]
		}
	}
}

func createCounter() {
	counterHTML = ""
	for appNR := range availableApps {
		counterHTML = counterHTML + "<label>" + availableApps[appNR] + " wurde " + strconv.Itoa(restartCounter[appNR]) + " mal automatisch neu gestartet</label><br><br>"
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

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(stdout)
	cmd.Start()
	availableAppsNumberForCmd = append(availableAppsNumberForCmd, procStart)
	line, _, _ := reader.ReadLine()
	stdOutput = append(stdOutput, string(line))

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
		writer.WriteString(stdins[availableAppsNumberForCmd[procToKill]] + "\n")
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

func procAppKillHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	procNr := q.Get("procNr")
	if procNr == "" {
		procNr = "notFound"
	}
	procToKill, err := strconv.Atoi(procNr)
	if err != nil {
		fmt.Println(err)
	} else {
		cmd := exec.Command(stopCommands[availableAppsNumberForCmd[procToKill]])
		cmd.Start() //es wird angenommen, dass die Beenden-Anwendung sich selbst schließt
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

	createCounter()
	createStopButtons()
	createShellOutputButtons()
	responseString := responseStringFirstLine +
		availableAppsButtonHTML +
		counterHTML +
		stopProcessButtons +
		outputButtonHTML +
		"<label>Letzter bekannter Output: " + stdOutput[procToOutput] + "</label><br><br>" +
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter XML Path: ")
	line, _ := reader.ReadString('\n')
	line = strings.Trim(line, "\r")
	line = strings.Replace(line, "\r", "", -1)
	line = strings.Replace(line, "\n", "", -1)
	path = line
	readXML(true, path)

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/procKill/", procKillHandler)
	http.HandleFunc("/procSoftKill/", procSoftKillHandler)
	http.HandleFunc("/procAppKill/", procAppKillHandler)
	http.HandleFunc("/procStart/", procStartHandler)
	http.HandleFunc("/procShellOutput/", procShellOutputHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))

}
