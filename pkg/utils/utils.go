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

// reverse reverses a slice of any type
func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}
