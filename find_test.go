package pdfrenderer_test

import (
	"os"
	"path/filepath"
	"testing"

	pdfrenderer "github.com/Open-Event-Systems/chromium-pdf-renderer"
)

func TestFind(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("PATH", tmpDir)

	t.Run("find exec name", func(t *testing.T) {
		fn := filepath.Join(tmpDir, "chrome")
		f, err := os.OpenFile(fn, os.O_CREATE|os.O_RDWR, 0o755)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fn)
		f.Close()

		res := pdfrenderer.FindChromium()
		if res != fn {
			t.Fatalf("expected %s, got %s", fn, res)
		}
	})

	t.Run("find nothing", func(t *testing.T) {
		res := pdfrenderer.FindChromium()
		if res != "" {
			t.Fatalf("expected empty string, got %s", res)
		}
	})

	t.Run("find exec in order", func(t *testing.T) {
		fn := filepath.Join(tmpDir, "chromium")
		f, err := os.OpenFile(fn, os.O_CREATE|os.O_RDWR, 0o755)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fn)
		f.Close()

		fn = filepath.Join(tmpDir, "chromium-browser")
		f, err = os.OpenFile(fn, os.O_CREATE|os.O_RDWR, 0o755)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fn)
		f.Close()

		res := pdfrenderer.FindChromium()
		if res != fn {
			t.Fatalf("expected %s, got %s", fn, res)
		}
	})
}
