package typst

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var sampleSource []byte

func init() {
	var err error
	sampleSource, err = os.ReadFile("testdata/sample.typ")
	if err != nil {
		panic("failed to read testdata/sample.typ: " + err.Error())
	}
}

func BenchmarkCompileBytes_Library(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}
	defer c.Close()

	doc, err := c.CompileBytes(sampleSource)
	if err != nil {
		b.Fatalf("warmup failed: %v", err)
	}
	b.Logf("PDF size: %d bytes", doc.Len())
	doc.Close()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		doc, err := c.CompileBytes(sampleSource)
		if err != nil {
			b.Fatalf("compile failed: %v", err)
		}
		doc.Close()
	}
}

func BenchmarkCompile_CLI(b *testing.B) {
	typstBin, err := exec.LookPath("typst")
	if err != nil {
		b.Skip("typst CLI not found")
	}

	srcPath, _ := filepath.Abs("testdata/sample.typ")
	outPath := filepath.Join(b.TempDir(), "output.pdf")

	cmd := exec.Command(typstBin, "compile", srcPath, outPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		b.Fatalf("warmup failed: %v\n%s", err, out)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		cmd := exec.Command(typstBin, "compile", srcPath, outPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			b.Fatalf("compile failed: %v\n%s", err, out)
		}
	}
}
