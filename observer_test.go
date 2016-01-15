package main

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
	"time"
)

//Testfall, ob xml richtig ausgelesen wird
func TestReadXML(t *testing.T) {
	//   !!!!!!!   EINTRAGEN DES NEUEN PFADES DER XML-DATEI   !!!!!!!
	readXML(true, "C:\\Go\\src\\_Programme\\Project\\observer.xml")
	//Test erfolgreich, falls Liste der availableApps mit Anwendungen gefüllt wurde
	if availableApps == nil || availableAppsButtonHTML == "" {
		t.Error("Fehler beim Einlesen der XML-Datei")
	}
}

//Testfall, ob Counter richtig funktioniert
func TestCreateCounter(t *testing.T) {
	createCounter()
	//Test erfolgreich, falls Counter angelegt wurde
	if counterHTML == "" {
		t.Error("Fehler beim Erzeugen eines Counters")
	}
}

//Testfall, ob MainHandler funktioniert
func TestMainHandler(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(mainHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	//Test erfolgreich, falls kein Error zurückgegeben wird
	if err != nil {
		t.Error("Fehler beim MainHandler!")
	}
}

//Testfall, ob ProcStartHandler ohne übergebene Programm-ID richtig funktioniert
func TestProcStartHandlerEmpty(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	readXML(false, "")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procStartHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	//Test erfolgreich, falls kein Fehler ausgegeben
	//--> erwartete Ausgabe des Programms = Fehler bei strconv
	if err != nil {
		t.Error("Fehler beim ProcStartHandlerEmpty!")
	}
}

//Testfall, ob ProcStarthandler mit Programm-ID und Automatischem Neustart richtig funktioniert
func TestProcStartHandlerWithAutoNewstart(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	readXML(false, "")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procStartHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	//Programm-ID wird übergeben und Neustart wird auf true gesetzt
	url := testServer.URL + "/?procStartID=0&autoRestart=true"
	_, err := http.Get(url)
	//Test erfolgreich, falls keine Fehler ausgegeben
	if err != nil {
		t.Error("Fehler beim ProcStartHandlerWithAutoNewstart!")
	}
}

//Testfall, ob ProcKillHandler mit Programm-ID richtig funktioniert
func TestProcKillHandlerWithProcNr(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	readXML(false, "")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung dees Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	//Programm-ID wird übergeben
	url := testServer.URL + "/?procNr=0"
	//Test erfolgereich, falls keine Fehler ausgegeben
	_, err := http.Get(url)
	if err != nil {
		t.Error("Fehler beim ProcKillHandlerWithProcNr!")
	}
}

//Testfall, ob ProcKillHandler ohne übergebene Programm-ID richtig funktioniert
func TestProcKillHandlerEmptyProcNr(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	readXML(false, "")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	//Test erfolgreich, falls keine Fehler ausgegeben
	//--> erwartete Ausgabe des Programms: Fehler bei strconv
	if err != nil {
		t.Error("Fehler beim ProcKillHandlerEmptyProcNr!")
	}
}

//Testfall, ob ProcSoftKillHandler ohne übergebene Programm-ID richtig funktioniert
func TestProcSoftKillHandlerEmptyProcNr(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procSoftKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	//Test erfolgreich, falls keine Fehler ausgegeben
	if err != nil {
		t.Error("Fehler beim ProcSoftKillHandlerEmptyProcNr!")
	}
}

//Testfall, ob ProcSoftKillHandler mit Programm-ID richtig funktioniert
func TestProcSoftKillHandlerWithProcNr(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	readXML(false, "")
	//zuerst wird ProcStartHandler aufgerufen, um ein Programm zu starten, welches anschließend beendbar ist
	help := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procStartHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServerHelp := httptest.NewServer(help.Handler)
	//Programm-ID wird übergeben und automatischer Neustart auf true gesetzt
	urlHelp := testServerHelp.URL + "/?procStartID=0&autoRestart=true"
	_, errHelp := http.Get(urlHelp)
	//Test läuft weiter, falls keine Fehlermeldung ausgegeben
	if errHelp != nil {
		t.Error("Fehler beim ProcStartHandlerWithAutoNewstart!")
	}
	//Aufruf des ProcSoftKillHandler
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procSoftKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	//Programm-Id wird übergeben
	url := testServer.URL + "/?procNr=0"
	_, err := http.Get(url)
	//Test erfolgreich, falls keine Fehlermeldung ausgegeben
	if err != nil {
		t.Error("Fehler beim ProcSoftKillHandlerWithProcNr!")
	}
}

//Testfall, ob ProcAppKillHandler mit Programm-ID richtig funktioniert
func TestProcAppKillHandlerWithProcNr(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	readXML(false, "")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procAppKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	//Programm-ID wird übergeben
	url := testServer.URL + "/?procNr=0"
	_, err := http.Get(url)
	//Test erfolgreich, falls keine Fehlermeldung ausgegeben
	if err != nil {
		t.Error("Fehler beim ProcAppKillHandlerWithProcNr!")
	}
}

//Testfall, ob ProcAppKillHandler ohne übergebene Programm-ID richtig funktioniert
func TestProcAppKillHandlerEmptyProcNr(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	readXML(false, "")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procAppKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	//Testfall erfolgreich, falls keine Fehlermeldung ausgegeben
	//--> erwartete Ausgabe des Programms: Fehler bei strconv
	if err != nil {
		t.Error("Fehler beim ProcKillHandlerEmptyProcNr!")
	}
}

//Testfall, ob ProcShellOutput mit Programm-ID richtig funktioniert
func TestProcShellOutputWithProcNr(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	readXML(false, "")
	//zuerst wird ProcStartHandler aufgerufen, um ein Programm zu starten, welches einen Output hat
	help := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procStartHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServerHelp := httptest.NewServer(help.Handler
	urlHelp := testServerHelp.URL + "/?procStartID=0&autoRestart=true"
	_, errHelp := http.Get(urlHelp)
	//Test läuft weiter, falls keine Fehlermeldung ausgegeben
	if errHelp != nil {
		t.Error("Fehler beim ProcStartHandlerWithAutoNewstart!")
	}
	//Aufruf des ProcShellOutputHandler
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procShellOutputHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	//Programm-ID wird übergeben
	url := testServer.URL + "/?procNr=0"
	_, err := http.Get(url)
	//Test erfolgreich, falls keine Fehlermeldung ausgegeben
	if err != nil {
		t.Error("Fehler beim ProcShellOutputHandlerWithProcNr!")
	}
}

//Testfall, ob ShellOutputButtons richtig erzeugt werden
func TestCreateShellOutputButton(t *testing.T) {
	//Prozesse werden in Liste angelegt
	runningProcesses = append(runningProcesses, exec.Command(availableApps[0]))
	createShellOutputButtons()
	//Testfall erfolgreich, falls outputButtonHTML gefüllt wurde
	if outputButtonHTML == "" {
		t.Error("Fehler beim Erzeugen eines Buttons für den Shell-Output")
	}
}
