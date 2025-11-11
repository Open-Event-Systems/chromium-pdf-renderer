package main

import (
	"flag"
	"io"
	"log"
	"os"
	"pdfrenderer"
)

func main() {
	var inputVal string
	var outputPath string

	flag.StringVar(&inputVal, "input", "", "the input path or URL")
	flag.StringVar(&outputPath, "output", "output.pdf", "the output path")

	flag.Parse()

	if inputVal == "" || outputPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	execPath := os.Getenv("CHROMIUM")
	if execPath == "" {
		execPath = pdfrenderer.FindChromium()
	}

	if execPath == "" {
		log.Println("could not find chromium executable")
		os.Exit(1)
	}

	renderer := pdfrenderer.NewRenderer(
		execPath,
		pdfrenderer.WithVirtualTimeBudget(10000),
	)

	_, stderr, err := renderer.Render(inputVal, outputPath, "")
	if err != nil {
		log.Printf("chromium exited with error: %s", err)
		log.Println("stderr:")
		stderrData, _ := io.ReadAll(stderr)
		os.Stderr.Write(stderrData)
		os.Exit(1)
	}
}
