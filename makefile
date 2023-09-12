.PHONY: all windows linux mac mac-arm clean

# Default target to build all architectures
all: windows linux mac mac-arm

# Build for Windows
windows:
	GOOS=windows GOARCH=amd64 go build -o ./bin/barcode_windows_amd64.exe main.go

# Build for Linux
linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/barcode_linux_amd64 main.go

# Build for macOS (Intel)
mac:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/barcode_darwin_amd64 main.go

# Build for macOS (ARM)
mac-arm:
	GOOS=darwin GOARCH=arm64 go build -o ./bin/barcode_darwin_arm64 main.go

# Clean up
clean:
	@rm -rf ./bin/*