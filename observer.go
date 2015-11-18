package observer

import (
	"fmt"
	"os/exec"
	"strings"
)

func openShell(list string) bool {
	cmd := exec.Command(list)
	stdout, err := cmd.Output()
	if err != nil {
		//fmt.Println(err)
		return false
	}
	fmt.Println(string(stdout))
	return true
}

func enterCommand(command string, location string) bool {
	cmd := exec.Command(command, location)

	cmd.Run()
	stdout, err := cmd.Output()
	if !strings.Contains(err.Error(), "asdf") {
		//fmt.Println(err)
		return false
	}
	fmt.Println(string(stdout))
	return true
}
