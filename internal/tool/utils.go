package tool

import (
	"fmt"
	"os/exec"
	"runtime"
)

// openURL opens the given URL or file with the system's default application
func openURL(urlOrPath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", urlOrPath)
	case "linux":
		cmd = exec.Command("xdg-open", urlOrPath)
	case "windows":
		cmd = exec.Command("cmd", "/c", fmt.Sprintf("start %q", urlOrPath))
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}
