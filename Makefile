.PHONY: build test clean rust bench

# Build the Rust static library (required before go build/test)
rust:
	cd typst-ffi && cargo build --release

# Run Go tests (builds Rust first)
test: rust
	go test -v -count=1 ./...

# Run benchmarks
bench: rust
	go test -bench=. -benchtime=10s -count=1 -run='^$$' -timeout=30m ./internal/customfont_test/

# Clean all build artifacts
clean:
	cd typst-ffi && cargo clean
	go clean -testcache

# Build everything
build: rust
	go build ./...
