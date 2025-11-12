# Chronapse


**Capture time. Compress moments. Command the lapse.**

Chronapse is a Linux-only timelapse recorder that blends the elegance of Go's [Bubbletea](https://github.com/charmbracelet/bubbletea) TUI with the simplicity of Python and OpenCV.
It lets you record long-duration timelapses directly from your laptop camera — efficiently, safely, and beautifully.

## Features

- Beautiful terminal UI with real-time progress tracking
- Efficient webcam capture using OpenCV
- Automatic video compilation with FFmpeg
- Configurable intervals and durations

## Todo
- [x] basic app
- [ ] live feed of recording
- [ ] improve TUI

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

## Requirements

- Linux (tested on Ubuntu/Debian)
- Webcam at `/dev/video0`
- Python 3.7+
- Go 1.21+
- FFmpeg
- OpenCV (opencv-python)

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions welcome! Please open issues or pull requests.
See [CONTRIBUTING.md](CONTRIBUTING.md) for more info.