package pdfrenderer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"slices"
)

// An object that allows rendering web content to PDF using chromium headless mode.
type Renderer struct {
	// The path to the chromium executable.
	Path string

	timeout           int
	virtualTimeBudget int
	logLevel          int

	runFunc func(*exec.Cmd) error
}

type rendererOption func(*Renderer)

var defaultArgs = []string{
	"--headless=new",
	"--disable-gpu",
	"--disable-dev-shm-usage",
	"--no-sandbox",
	"--no-pdf-header-footer",
}

// Create a new PDF renderer with the given options.
func NewRenderer(chromiumPath string, options ...rendererOption) *Renderer {
	renderer := &Renderer{
		Path:     chromiumPath,
		logLevel: 3,
	}

	for _, opt := range options {
		opt(renderer)
	}

	return renderer
}

// Set a timeout before the page is captured.
func WithTimeout(ms int) rendererOption {
	return func(r *Renderer) {
		r.timeout = ms
	}
}

// Set the virtual time budget before the page is captured.
func WithVirtualTimeBudget(ms int) rendererOption {
	return func(r *Renderer) {
		r.virtualTimeBudget = ms
	}
}

// Set the log level.
//
// 0 = INFO, 1 = WARNING, 2 = ERROR, 3 = FATAL
// The default log level is 1.
func WithLogLevel(level int) rendererOption {
	if level > 3 {
		level = 3
	} else if level < 0 {
		level = 0
	}

	return func(r *Renderer) {
		r.logLevel = level
	}
}

func (r *Renderer) getArgs(inputPath string, outputPath string) []string {
	args := slices.Clone(defaultArgs)

	args = append(args, fmt.Sprintf("--log-level=%d", r.logLevel))

	if r.timeout > 0 {
		args = append(args, fmt.Sprintf("--timeout=%d", r.timeout))
	}

	if r.virtualTimeBudget > 0 {
		args = append(args, fmt.Sprintf("--virtual-time-budget=%d", r.virtualTimeBudget))
	}

	// TODO: check on arg behavior on windows (do spaces etc. need escaped?)
	args = append(args, fmt.Sprintf("--print-to-pdf=%s", outputPath))
	args = append(args, inputPath)

	return args
}

// Render a page to PDF.
// input may be a file path or URL. outputPath must be a path to a file to be created and written to.
// Returns the stdout, stderr, and error.
// It is possible for this function to fail to create the output file without returning any error.
func (r *Renderer) RenderContext(ctx context.Context, input string, outputPath string, workDir string) (io.Reader, io.Reader, error) {
	args := r.getArgs(input, outputPath)
	cmd := exec.CommandContext(ctx, r.Path, args...)
	cmd.Dir = workDir

	stdoutBuf := bytes.NewBuffer(nil)
	stderrBuf := bytes.NewBuffer(nil)
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf

	runFunc := r.runFunc
	if r.runFunc == nil {
		runFunc = func(c *exec.Cmd) error {
			return c.Run()
		}
	}

	err := runFunc(cmd)
	return stdoutBuf, stderrBuf, err
}

// Render a page to PDF.
// input may be a file path or URL. outputPath must be a path to a file to be created and written to.
// Returns the stdout, stderr, and error.
// It is possible for this function to fail to create the output file without returning any error.
func (r *Renderer) Render(input string, outputPath string, workDir string) (io.Reader, io.Reader, error) {
	return r.RenderContext(context.Background(), input, outputPath, workDir)
}
