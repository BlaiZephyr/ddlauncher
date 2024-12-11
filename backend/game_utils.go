package backend

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func RunGameCommand() error {
	var execPath string

	switch runtime.GOOS {
	case "windows":
		execPath = filepath.Join(State.VersionsDir, fmt.Sprintf("DDNet-%s-win64", State.CurrentVersion), GameTitle+".exe")
	case "linux":
		execPath = filepath.Join(State.VersionsDir, fmt.Sprintf("DDNet-%s-linux_x86_64", State.CurrentVersion), GameTitle)
	default:
		return fmt.Errorf("unsupported OS")
	}

	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		_ = fmt.Errorf("game executable not found: %s", execPath)
		_, err := FetchDDNetZip(State.CurrentVersion)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(execPath)
	return cmd.Run()
}
