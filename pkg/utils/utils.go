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

func Ternary(b bool, x1 interface{}, x2 interface{}) interface{} {
	if b {
		return x1
	} else {
		return x2
	}
}
