package backend

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/google/go-github/v39/github"
)

func getCurrentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		_ = fmt.Errorf("damn.. %s", err)
	}

	versionsDir := filepath.Join(currentDir, "..", "Versions")
	absolutePath, _ := filepath.Abs(versionsDir)

	return absolutePath
}
func FetchDDNetZip(version string) (bool, error) {

	versionFolder := filepath.Join(getCurrentDir(), "Versions", fmt.Sprintf("DDNet-%s-%s", version, runtime.GOOS))

	if _, err := os.Stat(versionFolder); err == nil {
		executableSuffix := ".exe"
		if runtime.GOOS == "linux" {
			executableSuffix = ""
		}
		executablePath := filepath.Join(versionFolder, "DDNet"+executableSuffix)
		if _, err := os.Stat(executablePath); err == nil {
			fmt.Printf("Version %s already exists and appears to be fully downloaded. Skipping download.\n", version)
			return false, nil
		}
	}
	versionsDir := filepath.Join(getCurrentDir(), "Versions")
	if _, err := os.Stat(versionsDir); os.IsNotExist(err) {
		err := os.MkdirAll(versionsDir, 0755)
		if err != nil {
			return false, fmt.Errorf("failed to create Versions directory: %v", err)
		}
	}
	var downloadURL string
	switch os := runtime.GOOS; os {
	case "windows":
		downloadURL = fmt.Sprintf("https://ddnet.org/downloads/DDNet-%s-win64.zip", version)
	case "linux":
		downloadURL = fmt.Sprintf("https://ddnet.org/downloads/DDNet-%s-linux_x86_64.tar.xz", version)
	default:
		return false, fmt.Errorf("unsupported OS: %s", os)
	}

	fmt.Printf("Downloading %s\n", downloadURL)

	client := grab.NewClient()
	req, err := grab.NewRequest(getCurrentDir(), downloadURL)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}

	resp := client.Do(req)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("transferred %v / %v bytes (%.2f%%)\n", resp.BytesComplete(), resp.Size(), 100*resp.Progress())
		case <-resp.Done:
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		return false, fmt.Errorf("download failed: %v", err)
	}

	err = extractArchive(resp.Filename, getCurrentDir())
	if err != nil {
		return false, fmt.Errorf("extraction failed: %v", err)
	}

	return true, nil
}

func extractArchive(archivePath, destPath string) error {
	version := extractVersionFromFilename(archivePath)

	versionFolder := filepath.Join(destPath, fmt.Sprintf("DDNet-%s", version))
	err := os.MkdirAll(versionFolder, 0755)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-command", "Expand-Archive", "-Path", archivePath, "-DestinationPath", getCurrentDir(), "-Force")
	case "linux":
		cmd = exec.Command("tar", "-xvf", archivePath, "-C", versionFolder)
	default:
		return fmt.Errorf("unsupported OS")
	}

	output, err := cmd.CombinedOutput()
	fmt.Printf("Extraction output:\n%s\n", string(output))
	if err != nil {
		return fmt.Errorf("extraction failed: %v - output: %s", err, string(output))
	}

	if files, err := os.ReadDir(versionFolder); err == nil && len(files) == 0 {
		err = os.Remove(versionFolder)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
	err = os.Remove(archivePath)
	if err != nil {
		fmt.Printf(" %v\n", err)
	}

	return nil
}

func extractVersionFromFilename(filename string) string {
	parts := strings.Split(filepath.Base(filename), "-")
	if len(parts) >= 2 {
		return parts[1]
	}
	return "unknown"
}

func FetchGitHubTags() ([]string, error) {
	ctx := context.Background()
	client := github.NewClient(nil)

	opts := &github.ListOptions{
		Page:    1,
		PerPage: 50,
	}

	tags, resp, err := client.Repositories.ListTags(ctx, "ddnet", "ddnet", opts)
	if err != nil {
		fmt.Printf("Response Status: %v\n", resp.Status)
		return nil, err
	}

	var filteredTags []string

	for _, tag := range tags {
		if tag.Name != nil && strings.HasPrefix(*tag.Name, "1") {
			filteredTags = append(filteredTags, *tag.Name)
		}
	}

	return filteredTags, nil
}
