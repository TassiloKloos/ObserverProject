package main

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
	"time"
)

//	availableApps = append(availableApps, "C:\\Go\\src\\_Programme\\Project\\test.bat")
//	availableApps = append(availableApps, "C:\\Go\\src\\_Programme\\Project\\test2.bat")

func TestReadXML(t *testing.T) {
	readXML()
	if availableApps == nil || availableAppsButtonHTML == "" {
		t.Error("Fehler beim Einlesen der XML-Datei")
	}
}

func TestCreateStopButtons(t *testing.T) {
	runningProcesses = append(runningProcesses, exec.Command(availableApps[0]))
	createStopButtons()
	if stopProcessButtons == "" {
		t.Error("Fehler beim Erzeugen der Stop-Buttons")
	}
}
func TestCreateShellOutputButton(t *testing.T) {
	runningProcesses = append(runningProcesses, exec.Command(availableApps[0]))
	createShellOutputButtons()
	if outputButtonHTML == "" {
		t.Error("Fehler beim Erzeugen eines Buttons für den Shell-Output")
	}
}

func TestCreateCounter(t *testing.T) {
	createCounter()
	if counterHTML == "" {
		t.Error("Fehler beim Erzeugen eines Counters")
	}
}

func TestMainHandler(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(mainHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	if err != nil {
		t.Error("Fehler beim MainHandler!")
	}
}

// noch zu implementieren: if stdin, err := cmd.StdinPipe(); err != nil { 	Zeile 66
func TestProcStartHandlerEmpty(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procStartHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	if err != nil {
		t.Error("Fehler beim ProcStartHandlerEmpty!")
	}
}

func TestProcStartHandlerWithAutoNewstart(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procStartHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	url := testServer.URL + "/?procStartID=0&autoRestart=true"
	_, err := http.Get(url)
	if err != nil {
		t.Error("Fehler beim ProcStartHandlerWithAutoNewstart!")
	}
}

func TestProcKillHandlerEmptyProcNr(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	if err != nil {
		t.Error("Fehler beim ProcKillHandlerEmptyProcNr!")
	}
}

func TestProcKillHandlerWithProcNr(t *testing.T) { //Fehler
	runningProcesses = append(runningProcesses, exec.Command(availableApps[0]))
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	url := testServer.URL + "/?procNr=0"
	_, err := http.Get(url) //Fehler, obwohl Prozess vorhanden!!! --> runningProcesses[procToKill].Process.Kill(), Zeile 132
	if err != nil {
		t.Error("Fehler beim ProcKillHandlerWithProcNr!")
	}
}

func TestProcSoftKillHandlerEmptyProcNr(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procSoftKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	if err != nil {
		t.Error("Fehler beim ProcSoftKillHandlerEmptyProcNr!")
	}
}

func TestProcSoftKillHandlerWithProcNr(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procSoftKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	url := testServer.URL + "/?procNr=0"
	_, err := http.Get(url) //Fehler --> ShellOutput[procToOutput], Zeile 197
	if err != nil {
		t.Error("Fehler beim ProcSoftKillHandlerWithProcNr!")
	}
}

func TestProcShellOutputWithProcNr(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procShellOutputHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	url := testServer.URL + "/?procNr=0"
	_, err := http.Get(url)
	if err != nil {
		t.Error("Fehler beim ProcShellOutputHandlerWithProcNr!")
	}
}
