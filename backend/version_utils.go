package backend

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/google/go-github/v39/github"
)

func FetchDDNetZip(version string) (bool, error) {

	versionFolder := filepath.Join(State.VersionsDir, "Versions", fmt.Sprintf("DDNet-%s-%s", version, runtime.GOOS))

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
	versionsDir := filepath.Join(State.VersionsDir, "Versions")
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
	req, err := grab.NewRequest(State.VersionsDir, downloadURL)
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

	err = extractArchive(resp.Filename)
	if err != nil {
		return false, fmt.Errorf("extraction failed: %v", err)
	}

	return true, nil
}

func extractArchive(archivePath string) error {

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-command", "Expand-Archive", "-Path", archivePath, "-DestinationPath", State.VersionsDir, "-Force")
	case "linux":
		cmd = exec.Command("tar", "-xvf", archivePath, "-C", State.VersionsDir)
	default:
		return fmt.Errorf("unsupported OS")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("extraction failed: %v - output: %s", err, string(output))
	}

	if files, err := os.ReadDir(State.VersionsDir); err == nil && len(files) == 0 {
		err = os.Remove(State.VersionsDir)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
	err = os.Remove(archivePath)
	if err != nil {
		fmt.Printf(" %v\n", err)
	}

	err = os.Remove(State.VersionsDir + "/Versions")
	if err != nil {
		fmt.Printf(" %v\n", err)
	}

	return nil
}

func filterAndSortTags(tags []*github.RepositoryTag) []string {
	var filteredTags []string

	fullVersionRegex := regexp.MustCompile(`^(\d+\.\d+(?:\.\d+)?)$`)

	for _, tag := range tags {
		if tag.Name == nil {
			continue
		}
		matches := fullVersionRegex.FindStringSubmatch(*tag.Name)
		if len(matches) > 0 {
			filteredTags = append(filteredTags, matches[1])
		}
	}

	sort.Slice(filteredTags, func(i, j int) bool {
		return compareVersions(filteredTags[i], filteredTags[j])
	})

	return filteredTags
}

func compareVersions(v1, v2 string) bool {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		num1 := 0
		num2 := 0

		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &num1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &num2)
		}

		if num1 > num2 {
			return true
		} else if num1 < num2 {
			return false
		}
	}

	return len(parts1) > len(parts2)
}

// TODO: this still blocks main loop.
func FetchGitHubTagsAsync(callback func([]string, error)) {
	go func() {
		tags, err := FetchGitHubTags()
		callback(tags, err)
	}()
}

func FetchGitHubTags() ([]string, error) {
	ctx := context.Background()
	client := github.NewClient(nil)

	opts := &github.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	tags, resp, err := client.Repositories.ListTags(ctx, "ddnet", "ddnet", opts)
	if err != nil {
		log.Printf("GitHub API Response Status: %v\n", resp.Status)
		return nil, err
	}

	filteredTags := filterAndSortTags(tags)

	log.Printf("Filtered Tags: %v\n", filteredTags)

	if len(filteredTags) == 0 {
		return nil, fmt.Errorf("no full version tags found")
	}

	return filteredTags, nil
}
