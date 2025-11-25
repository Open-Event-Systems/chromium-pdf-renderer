package pdfrenderer

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var execNames = []string{
	"chromium-browser",
	"google-chrome",
	"chromium",
	"chrome",
}

// Find the chromium/chrome executable.
// Returns the empty string if not found.
func FindChromium() string {
	for _, name := range execNames {
		path := findExec(name)
		if path != "" {
			return path
		}
	}

	if runtime.GOOS == "windows" {
		for _, path := range getWindowsPaths() {
			_, err := os.Stat(path)
			if err == nil {
				return path
			}
		}
	}

	return ""
}

func findExec(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return path
}

func getWindowsPaths() []string {
	var paths []string

	paths = append(paths, filepath.Join(os.Getenv("ProgramFiles"), "Google", "Chrome", "Application", "chrome.exe"))
	paths = append(paths, filepath.Join(os.Getenv("ProgramFiles(x86)"), "Google", "Chrome", "Application", "chrome.exe"))
	paths = append(paths, filepath.Join(os.Getenv("LocalAppData"), "Google", "Chrome", "Application", "chrome.exe"))

	return paths
}
