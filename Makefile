.PHONY: build run install clean test help

# Default target
help:
	@echo "Chronapse Timelapse Recorder - Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make build      - Build the Go TUI executable"
	@echo "  make run        - Build and run the application"
	@echo "  make install    - Install all dependencies"
	@echo "  make clean      - Remove build artifacts and temporary files"
	@echo "  make test       - Run a quick test recording (30 seconds)"
	@echo "  make help       - Show this help message"

# Build the Go executable
build:
	@echo "Building Chronapse..."
	go build -o chronapse main.go
	@echo "Build complete! Run with: ./chronapse"

# Build and run
run: build
	./chronapse

# Install all dependencies
install:
	@echo "Installing system dependencies..."
	@echo "Note: This requires sudo privileges"
	sudo apt update
	sudo apt install -y python3 python3-pip ffmpeg golang-go
	@echo ""
	@echo "Installing Python dependencies..."
	pip3 install -r requirements.txt
	@echo ""
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo ""
	@echo "Installation complete!"

# Clean build artifacts and temporary files
clean:
	@echo "Cleaning up..."
	rm -f chronapse
	rm -rf frames/
	rm -f *.mp4
	@echo "Cleanup complete!"

# Run a quick test
test:
	@echo "Running 30-second test timelapse..."
	python3 timelapse.py -i 2 -d 30 -o test.mp4
	@echo ""
	@echo "Test complete! Play with: vlc test.mp4"
