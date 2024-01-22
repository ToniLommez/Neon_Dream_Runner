package neon

import (
	"os"
	"os/exec"
	"runtime"
)

func clearScreen() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls") // Para Windows
	} else {
		cmd = exec.Command("clear") // Para sistemas Unix-like
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
