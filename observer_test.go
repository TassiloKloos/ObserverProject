package main

import (
	"fmt"
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
	res, err := http.Get(testServer.URL)
	if err != nil {
		t.Error("Fehler beim MainHandler!")
	}
	fmt.Println("Res.Body: ", res.Body)
}
func TestProcStartHandler(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procStartHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	res, err := http.Get(testServer.URL)
	if err != nil {
		t.Error("Fehler beim ProcStartHandler!")
	}
	fmt.Println("Res.Body: ", res.Body)
}

func TestProcKillHandler(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(procKillHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	testServer := httptest.NewServer(s.Handler)
	res, err := http.Get(testServer.URL)
	if err != nil {
		t.Error("Fehler beim ProcKillHandler!")
	}
	fmt.Println("Res.Body: ", res.Body)
}
