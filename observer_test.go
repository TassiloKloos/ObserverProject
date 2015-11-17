package observer

import "testing"

func TestOpenShell(t *testing.T) {
	var solve bool = openShell("tasklist")
	if solve != true {
		t.Error("Fehler beim Öffnen der Shell")
	}
}

func TestOpenShellWrong(t *testing.T) {
	var ergeb bool = openShell("x")
	if ergeb != false {
		t.Error("Öffnen der falschen Shell möglich")
	}
}

func TestEnterCommand(t *testing.T) {
	var solve bool = enterCommand("explorer.exe", "C:\\Go\\src\\_Programme\\Project")
	if solve != true {
		t.Error("Fehler beim Eingeben von Commands")
	}
}

func TestEnterWrongCommand(t *testing.T) {
	var solve bool = enterCommand("x", "y")
	if solve != false {
		t.Error("Eingeben von falschen Commands möglich")
	}
}
