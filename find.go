package pdfrenderer

import "os/exec"

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

	return ""
}

func findExec(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return path
}
