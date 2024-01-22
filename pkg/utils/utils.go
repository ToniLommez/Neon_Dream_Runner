package utils

import (
	"os"
	"os/exec"
	"runtime"
)

func ClearScreen() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls") // Para Windows
	} else {
		cmd = exec.Command("clear") // Para sistemas Unix-like
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
