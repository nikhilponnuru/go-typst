package typst

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func newTestCompiler(t *testing.T) *Compiler {
	t.Helper()
	c, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	t.Cleanup(func() { c.Close() })
	return c
}

func TestCompile_simple(t *testing.T) {
	c := newTestCompiler(t)
	doc, err := c.Compile(strings.NewReader(`= Hello, Typst!

This is a simple document.
`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer doc.Close()

	if doc.Len() == 0 {
		t.Fatal("expected non-empty PDF output")
	}
	if !bytes.HasPrefix(doc.Bytes(), []byte("%PDF-")) {
		t.Fatalf("output does not look like a PDF, starts with: %q", doc.Bytes()[:min(20, doc.Len())])
	}
}

func TestCompile_emptySource(t *testing.T) {
	c := newTestCompiler(t)
	_, err := c.Compile(strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error for empty source")
	}
	var ce *CompileError
	if !asCompileError(err, &ce) {
		t.Fatalf("expected CompileError, got %T: %v", err, err)
	}
}

func TestCompile_invalidTypst(t *testing.T) {
	c := newTestCompiler(t)
	doc, err := c.Compile(strings.NewReader(`#let x = `))
	if err == nil {
		doc.Close()
	}
}

func TestCompileBytes(t *testing.T) {
	c := newTestCompiler(t)
	doc, err := c.CompileBytes([]byte("Hello from bytes!"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer doc.Close()

	if !bytes.HasPrefix(doc.Bytes(), []byte("%PDF-")) {
		t.Fatal("output does not look like a PDF")
	}
}

func TestDocument_Read(t *testing.T) {
	c := newTestCompiler(t)
	doc, err := c.CompileBytes([]byte("Hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer doc.Close()

	got, err := io.ReadAll(doc)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if !bytes.HasPrefix(got, []byte("%PDF-")) {
		t.Fatal("output does not look like a PDF")
	}
	if len(got) != doc.Len() {
		t.Fatalf("ReadAll returned %d bytes, expected %d", len(got), doc.Len())
	}
}

func TestDocument_WriteTo(t *testing.T) {
	c := newTestCompiler(t)
	doc, err := c.CompileBytes([]byte("Hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer doc.Close()

	var buf bytes.Buffer
	n, err := doc.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}
	if int(n) != doc.Len() {
		t.Fatalf("WriteTo wrote %d bytes, expected %d", n, doc.Len())
	}
	if !bytes.HasPrefix(buf.Bytes(), []byte("%PDF-")) {
		t.Fatal("output does not look like a PDF")
	}
}

func TestDocument_CloseIdempotent(t *testing.T) {
	c := newTestCompiler(t)
	doc, err := c.CompileBytes([]byte("Hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doc.Close()
	doc.Close() // should not panic

	if doc.Len() != 0 {
		t.Fatal("expected 0 length after close")
	}
	if doc.Bytes() != nil {
		t.Fatal("expected nil bytes after close")
	}
}

func TestDocument_ReadAfterClose(t *testing.T) {
	c := newTestCompiler(t)
	doc, err := c.CompileBytes([]byte("Hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	doc.Close()

	_, err = doc.Read(make([]byte, 10))
	if err == nil {
		t.Fatal("expected error reading closed document")
	}
}

func TestDocument_WriteTo_ioCopy(t *testing.T) {
	c := newTestCompiler(t)
	doc, err := c.Compile(strings.NewReader(`
= Report

#lorem(200)
`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer doc.Close()

	var buf bytes.Buffer
	n, err := io.Copy(&buf, doc)
	if err != nil {
		t.Fatalf("io.Copy error: %v", err)
	}
	if int(n) != doc.Len() {
		t.Fatalf("io.Copy wrote %d bytes, expected %d", n, doc.Len())
	}
}

func TestCompiler_CloseIdempotent(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}
	c.Close()
	c.Close() // should not panic
}

func TestCompiler_CompileAfterClose(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}
	c.Close()

	_, err = c.CompileBytes([]byte("Hello"))
	if err == nil {
		t.Fatal("expected error compiling with closed compiler")
	}
}

func TestMultipleCompilers(t *testing.T) {
	c1, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer c1.Close()

	c2, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer c2.Close()

	doc1, err := c1.CompileBytes([]byte("From compiler 1"))
	if err != nil {
		t.Fatal(err)
	}
	defer doc1.Close()

	doc2, err := c2.CompileBytes([]byte("From compiler 2"))
	if err != nil {
		t.Fatal(err)
	}
	defer doc2.Close()

	if !bytes.HasPrefix(doc1.Bytes(), []byte("%PDF-")) || !bytes.HasPrefix(doc2.Bytes(), []byte("%PDF-")) {
		t.Fatal("one or both outputs are not PDFs")
	}
}

func asCompileError(err error, target **CompileError) bool {
	if ce, ok := err.(*CompileError); ok {
		*target = ce
		return true
	}
	return false
}
