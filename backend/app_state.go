package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	WindowTitle = "DDLauncher"
	GameTitle   = "DDNet"
)

type AppState struct {
	LogfilePath     string
	GamePath        string
	LatestVersion   string
	SelectedVersion string
	CurrentVersion  string
	UserDir         string
	VersionsDir     string
}

var State AppState

func InitAppState() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("home directory be missing yo: %v\n", err)
		return
	}

	versionsDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("wd directory be missing yo: %v\n", err)
	} else {
		versionsDir += "/Versions"
	}

	State.VersionsDir = versionsDir
	State.UserDir = homeDir

	//fully unused atm.
	switch os := runtime.GOOS; os {
	case "windows":
		State.LogfilePath = filepath.Join(State.UserDir, "AppData", "Roaming", "DDNet", "logfile.txt")
	case "linux":
		State.LogfilePath = filepath.Join(".local", "share", "ddnet", "logfile.txt")
	default:
		fmt.Printf("Unsupported OS: %s\n", os)
	}
}
