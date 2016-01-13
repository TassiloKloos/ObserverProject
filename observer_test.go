package main

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
	"time"
)

func TestReadXML(t *testing.T) {
	readXML()
	if availableApps == nil || availableAppsButtonHTML == "" {
		t.Error("Fehler beim Einlesen der XML-Datei")
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
	readXML()
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
	readXML()
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

func TestProcKillHandlerWithProcNr(t *testing.T) { //Fehler
	readXML()
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	url := testServer.URL + "/?procNr=0"
	_, err := http.Get(url) //Fehler, Prozess nicht vorhanden!!! --> runningProcesses[procToKill].Process.Kill(), Zeile 163
	if err != nil {
		t.Error("Fehler beim ProcKillHandlerWithProcNr!")
	}
}

func TestProcKillHandlerEmptyProcNr(t *testing.T) {
	readXML()
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
	readXML()
	help := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procStartHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServerHelp := httptest.NewServer(help.Handler)
	urlHelp := testServerHelp.URL + "/?procStartID=0&autoRestart=true"
	_, errHelp := http.Get(urlHelp)
	if errHelp != nil {
		t.Error("Fehler beim ProcStartHandlerWithAutoNewstart!")
	}
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procSoftKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	url := testServer.URL + "/?procNr=0"
	_, err := http.Get(url)
	if err != nil {
		t.Error("Fehler beim ProcSoftKillHandlerWithProcNr!")
	}
}

func TestProcShellOutputWithProcNr(t *testing.T) { //Noch nicht implementierbar!!!
	readXML()
	help := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procStartHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServerHelp := httptest.NewServer(help.Handler)
	urlHelp := testServerHelp.URL + "/?procStartID=0&autoRestart=true"
	_, errHelp := http.Get(urlHelp)
	if errHelp != nil {
		t.Error("Fehler beim ProcStartHandlerWithAutoNewstart!")
	}
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

func TestCreateShellOutputButton(t *testing.T) {
	runningProcesses = append(runningProcesses, exec.Command(availableApps[0]))
	createShellOutputButtons()
	if outputButtonHTML == "" {
		t.Error("Fehler beim Erzeugen eines Buttons für den Shell-Output")
	}
}

func TestCreateStopButtons(t *testing.T) {
	readXML()
	runningProcesses = append(runningProcesses, exec.Command(availableApps[0]))
	//	createStopButtons() //Fehler bei createStopButtons(),
	if stopProcessButtons == "" {
		t.Error("Fehler beim Erzeugen der Stop-Buttons")
	}
}
