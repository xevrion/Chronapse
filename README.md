# Chronapse

**Capture time. Compress moments. Command the lapse.**

Chronapse is a Linux-only timelapse recorder that blends the elegance of Go's [Bubbletea](https://github.com/charmbracelet/bubbletea) TUI with the simplicity of Python and OpenCV.
It lets you record long-duration timelapses directly from your laptop camera — efficiently, safely, and beautifully.

## Features

- Beautiful terminal UI with real-time progress tracking
- Efficient webcam capture using OpenCV
- Automatic video compilation with FFmpeg
- Graceful shutdown and error handling
- Configurable intervals and durations
- Production-ready and idiomatic code

## Quick Start

### Prerequisites

```bash
# Install dependencies
sudo apt install python3 python3-pip ffmpeg golang-go

# Install Python packages
pip3 install -r requirements.txt
```

### Build and Run

```bash
# Build the TUI
go build -o chronapse main.go

# Run
./chronapse
```

### Usage

1. Set your recording interval (seconds between frames)
2. Set total duration (in seconds)
3. Choose output filename
4. Press Enter to start recording
5. Press 'q' to stop early (optional)

### Example

Record a 10-minute timelapse with 5-second intervals:
- Interval: `5`
- Duration: `600`
- Output: `sunset.mp4`

Result: 120 frames → 4-second video at 30fps

## Documentation

See [USAGE.md](USAGE.md) for:
- Complete architecture explanation
- Detailed usage examples
- Safety recommendations for long recordings
- Troubleshooting guide
- Performance optimization tips

## Project Structure

```
Chronapse/
├── main.go           # Go Bubbletea TUI (frontend)
├── timelapse.py      # Python recorder (backend)
├── requirements.txt  # Python dependencies
├── go.mod           # Go module definition
├── README.md        # This file
└── USAGE.md         # Comprehensive documentation
```

## How It Works

1. **Go TUI** collects user input and spawns Python subprocess
2. **Python backend** captures frames from webcam at specified intervals
3. **Progress updates** stream from Python to Go via stdout
4. **FFmpeg** compiles frames into MP4 video
5. **Cleanup** removes temporary frames automatically

## Requirements

- Linux (tested on Ubuntu/Debian)
- Webcam at `/dev/video0`
- Python 3.7+
- Go 1.21+
- FFmpeg
- OpenCV (opencv-python)

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions welcome! Please open issues or pull requests.