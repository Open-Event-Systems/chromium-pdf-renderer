package pdfrenderer

import (
	"bytes"
	"io"
	"os/exec"
	"slices"
	"testing"
)

func TestRender(t *testing.T) {
	r := NewRenderer(
		"/test",
		WithLogLevel(3),
		WithTimeout(1000),
		WithVirtualTimeBudget(5000),
	)

	var resCmd *exec.Cmd

	r.runFunc = func(c *exec.Cmd) error {
		resCmd = c
		c.Stdout.Write([]byte("stdout"))
		c.Stderr.Write([]byte("stderr"))
		return nil
	}

	stdout, stderr, err := r.RenderContext(t.Context(), "/input", "/output", "/workdir")
	if err != nil {
		t.Fatal(err)
	}

	if resCmd.Path != "/test" {
		t.Fatalf("expected /test, got %s", resCmd.Path)
	}

	if resCmd.Dir != "/workdir" {
		t.Fatalf("expected /workdir, got %s", resCmd.Dir)
	}

	args := slices.Clone(resCmd.Args[1:])
	slices.Sort(args)

	expectedArgs := []string{
		"--headless=new",
		"--disable-gpu",
		"--disable-dev-shm-usage",
		"--no-sandbox",
		"--no-pdf-header-footer",
		"--log-level=3",
		"--timeout=1000",
		"--virtual-time-budget=5000",
		"--print-to-pdf=/output",
		"/input",
	}
	slices.Sort(expectedArgs)

	if !slices.Equal(expectedArgs, args) {
		t.Fatalf("expected %s, got %s", expectedArgs, args)
	}

	stdoutVal, err := io.ReadAll(stdout)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte("stdout"), stdoutVal) {
		t.Fatalf("expected stdout, got %s", stdoutVal)
	}

	stderrVal, err := io.ReadAll(stderr)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte("stderr"), stderrVal) {
		t.Fatalf("expected stderr, got %s", stdoutVal)
	}
}