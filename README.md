# go-typst

A Go library for compiling [Typst](https://typst.app) markup into PDF, powered by the Typst Rust crate (v0.14) via CGO.

- **Zero-copy output** — PDF bytes stay in Rust memory, never copied to the Go heap.
- **Instance-based** — each `Compiler` has its own fonts and caches, safe for concurrent use.
- **Custom fonts** — load any TTF/OTF when creating a `Compiler`, alongside the bundled defaults.

## Prerequisites

- **Rust** toolchain (`rustc`, `cargo`) — [install via rustup](https://rustup.rs)
- **Go** 1.21+
- **C compiler** (for CGO linking)

## Build

Build the Rust static library first:

```bash
make rust
```

## Usage

### Basic

```go
package main

import (
	"fmt"
	"os"
	"strings"

	typst "github.com/sarat/go-typst"
)

func main() {
	// Create a compiler (bundled fonts only).
	c, err := typst.New()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// Compile from an io.Reader.
	doc, err := c.Compile(strings.NewReader(`
= My Document

Hello from *Typst* compiled via Go!

$ sum_(k=1)^n k = (n(n+1))/2 $
`))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer doc.Close()

	// Zero-copy write directly from Rust memory to file.
	f, _ := os.Create("output.pdf")
	defer f.Close()
	doc.WriteTo(f)

	fmt.Printf("Written %d bytes to output.pdf\n", doc.Len())
}
```

### Custom Fonts

```go
regular, _ := os.ReadFile("fonts/Inter.ttf")
italic, _ := os.ReadFile("fonts/Inter-Italic.ttf")

// Bundled fonts are always included; custom fonts are added on top.
c, err := typst.New(regular, italic)
if err != nil {
    panic(err)
}
defer c.Close()

doc, _ := c.CompileBytes([]byte(`#set text(font: "Inter"); Hello!`))
defer doc.Close()
```

### Multiple Independent Compilers

```go
// Each compiler has its own fonts and caches — no shared state.
invoiceCompiler, _ := typst.New(invoiceFonts...)
reportCompiler, _ := typst.New(reportFonts...)

// Safe to use concurrently from different goroutines.
go func() { invoiceCompiler.CompileBytes(invoiceSource) }()
go func() { reportCompiler.CompileBytes(reportSource) }()
```

## API

### `func New(fonts ...[]byte) (*Compiler, error)`

Creates a new independent compiler instance. Bundled fonts (Libertinus Serif, New Computer Modern, DejaVu Sans Mono) are always loaded. Any additional font bytes passed here (TTF/OTF) are loaded on top.

### `type Compiler`

```go
func (c *Compiler) Compile(r io.Reader) (*Document, error)    // from io.Reader
func (c *Compiler) CompileBytes(source []byte) (*Document, error) // from []byte (fastest)
func (c *Compiler) Close() error                               // frees all resources
```

- **`Compile(r)`** — reads all bytes from `r`, compiles to PDF.
- **`CompileBytes(b)`** — compiles directly from a byte slice. Fastest path — avoids `io.ReadAll`.
- **`Close()`** — frees the compiler and all its internal resources. Idempotent. A runtime finalizer acts as safety net.

A `Compiler` is safe for concurrent use from multiple goroutines.

### `type Document`

```go
func (d *Document) Bytes() []byte                       // zero-copy view into Rust memory
func (d *Document) Len() int                            // PDF size in bytes
func (d *Document) Read(p []byte) (int, error)          // io.Reader
func (d *Document) WriteTo(w io.Writer) (int64, error)  // io.WriterTo (zero-copy)
func (d *Document) Close() error                        // frees Rust memory
```

- **`Bytes()`** — returns a slice backed directly by Rust-allocated memory. No allocation, no copy. Valid until `Close()`.
- **`WriteTo(w)`** — writes the PDF directly from Rust memory to `w`. Fastest path for writing to a file — single write, no Go heap allocation.
- **`Read(p)`** — standard `io.Reader`. Works with `io.Copy` etc.
- **`Close()`** — frees the underlying Rust memory. Idempotent.

### `type CompileError`

```go
type CompileError struct {
    Message string
}
```

Returned when Typst compilation or PDF export fails.

## Memory Model

```
New(fonts...)
  └─ Rust: parses fonts, builds library → heap-allocated Compiler instance

c.CompileBytes(source)
  ├─ Rust: copies source, compiles → Rust-allocated PDF bytes
  └─ Returns *Document pointing directly at Rust memory (zero-copy)

doc.WriteTo(file)
  └─ Writes from Rust memory → fd (single write syscall, no Go allocation)

doc.Close()
  └─ Frees Rust-allocated PDF memory

c.Close()
  └─ Frees compiler (fonts, library, caches)
```

The only copy is the source input. The PDF output is never copied into Go heap memory.

## Benchmarks

Measured on Apple M3 Pro (11 cores), Typst 0.14.2, `go test -bench`:

### Serial (single goroutine)

| Document | Library | CLI | Speedup | Docs/sec |
|---|---|---|---|---|
| 1-page (custom font, 19 KB) | **0.25 ms** | 57 ms | **228×** | 3,960 |
| 5-page report (109 KB) | **2.22 ms** | 90 ms | **41×** | 450 |
| 1000-page table (94 MB) | **0.97 s** | 9.3 s | **9.6×** | 1.03 |

### Parallel (GOMAXPROCS=11)

| Document | Shared Compiler | Compiler-per-CPU | Docs/sec |
|---|---|---|---|
| 1-page (custom font) | 0.246 ms/op | 0.245 ms/op | **~44,900** |
| 5-page report | 8.50 ms/op | 8.50 ms/op | **~1,295** |
| 1000-page table | 0.64 s/op | 0.63 s/op | **~5.3** |

Go-side allocations: **1 alloc, 48 bytes** per compile (just the `Document` struct).

## Building & Distribution

Rust is only needed **once** to build the static library. The final Go binary is fully self-contained — no Rust runtime, no shared libraries, no external dependencies.

```bash
# 1. Build the Rust static library (one time, or in CI)
make rust

# 2. Build your Go binary as usual — the Typst compiler is statically linked in
go build -o myapp ./cmd/myapp

# The resulting binary is standalone. Ship it anywhere.
./myapp
```

To distribute to teammates or CI environments that don't have Rust:

1. Build `libtypst_ffi.a` once on a machine with Rust (or in CI).
2. Commit it to your project or store it as a build artifact.
3. Point the CGO linker at it — consumers only need Go and a C compiler.

```bash
# Example: copy the prebuilt library into your project
cp path/to/go-typst/typst-ffi/target/release/libtypst_ffi.a ./vendor-lib/

# Build with custom library path
CGO_LDFLAGS="-L./vendor-lib -ltypst_ffi -lm -liconv -framework CoreFoundation -framework Security" \
    go build ./...
```

For cross-compilation, build the Rust library for each target:

```bash
# Example: build for linux/amd64 from macOS
rustup target add x86_64-unknown-linux-gnu
cd typst-ffi && cargo build --release --target x86_64-unknown-linux-gnu
# Output: typst-ffi/target/x86_64-unknown-linux-gnu/release/libtypst_ffi.a
```

## Testing

```bash
make test
```

## Limitations

- **Single-file only**: no support for `#import` of other Typst files or packages.
- **No external files**: the `file()` World method returns "not found" — embedded images from external paths won't work.
- **Fonts fixed at creation**: fonts are loaded once when the `Compiler` is created via `New()` and cannot be changed afterwards.
- **macOS/Linux**: tested on macOS (arm64). Linux should work but may need adjustments to linker flags.
