package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestReadXML(t *testing.T) {
	readXML()
	if availableApps == nil || availableAppsButtonHTML == "" {
		t.Error("Fehler beim Einlesen der XML-Datei")
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

func TestProcKillHandlerWithProcNr(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	url := testServer.URL + "/?procNr=0"
	_, err := http.Get(url)
	if err != nil {
		t.Error("Fehler beim ProcKillHandlerWithProcNr!")
	}
}

// Main ist nicht testbar!! --> Auslagern alles unnötigen Codes
