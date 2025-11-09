# Chronapse - Timelapse Recorder

A lightweight, efficient timelapse recorder for Linux systems with a beautiful terminal UI.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [How It Works](#how-it-works)
- [Safety Recommendations](#safety-recommendations)
- [Troubleshooting](#troubleshooting)

## Architecture Overview

Chronapse uses a two-part architecture:

```
┌─────────────────┐
│   Go TUI (UI)   │  - User input collection
│   (main.go)     │  - Real-time status display
│                 │  - Process management
└────────┬────────┘
         │ spawns subprocess
         ▼
┌─────────────────┐
│ Python Backend  │  - Webcam frame capture (OpenCV)
│ (timelapse.py)  │  - Video compilation (FFmpeg)
│                 │  - Progress reporting
└─────────────────┘
```

### Component Responsibilities

**Go TUI (`main.go`)**:
- Interactive menu for configuration (interval, duration, output file)
- Spawns Python subprocess with appropriate arguments
- Parses stdout/stderr for progress updates
- Displays real-time recording status
- Handles user interrupts (graceful shutdown)

**Python Backend (`timelapse.py`)**:
- Captures frames from `/dev/video0` using OpenCV
- Saves frames to temporary `frames/` directory
- Reports progress via stdout (parseable format)
- Compiles frames into MP4 using FFmpeg
- Cleans up temporary files
- Handles signals for graceful shutdown

## Prerequisites

### System Requirements

- **OS**: Linux (tested on Ubuntu/Debian)
- **Webcam**: `/dev/video0` or similar V4L2 device
- **Disk Space**: ~100MB minimum (varies with duration)

### Software Dependencies

#### Python 3
```bash
sudo apt update
sudo apt install python3 python3-pip
```

#### OpenCV
```bash
pip3 install opencv-python numpy
# OR using requirements.txt
pip3 install -r requirements.txt
```

#### FFmpeg
```bash
sudo apt install ffmpeg
```

#### Go (1.21+)
```bash
# If not already installed
sudo apt install golang-go
```

### Verify Installation

```bash
# Check Python
python3 --version

# Check OpenCV
python3 -c "import cv2; print(cv2.__version__)"

# Check FFmpeg
ffmpeg -version

# Check Go
go version

# Check webcam
ls -l /dev/video*
```

## Installation

1. **Clone the repository**:
```bash
cd ~/Coding
git clone <repository-url> Chronapse
cd Chronapse
```

2. **Install Python dependencies**:
```bash
pip3 install -r requirements.txt
```

3. **Install Go dependencies**:
```bash
go mod tidy
```

4. **Build the Go TUI**:
```bash
go build -o chronapse main.go
```

5. **Make Python script executable** (optional):
```bash
chmod +x timelapse.py
```

## Usage

### Quick Start

1. **Run the TUI**:
```bash
./chronapse
```

2. **Configure recording**:
   - **Interval**: Seconds between each frame (e.g., `5` = one frame every 5 seconds)
   - **Duration**: Total recording time in seconds (e.g., `600` = 10 minutes)
   - **Output file**: Video filename (e.g., `timelapse.mp4`)

3. **Start recording**:
   - Navigate with `Tab` or arrow keys
   - Press `Enter` on "Start Recording"

4. **Monitor progress**:
   - Watch the progress bar and frame counter
   - View recent activity logs
   - Press `q` to stop early (frames captured so far will still be compiled)

5. **Output**:
   - Video saved to specified path
   - Temporary frames automatically deleted

### Example Configurations

**Quick test (30 seconds, 2-second intervals)**:
- Interval: `2`
- Duration: `30`
- Output: `test.mp4`
- Result: 15 frames → 0.5-second video at 30fps

**One hour sunset (5-minute intervals)**:
- Interval: `300` (5 minutes)
- Duration: `3600` (1 hour)
- Output: `sunset.mp4`
- Result: 12 frames → 0.4-second video at 30fps

**8-hour workday (30-second intervals)**:
- Interval: `30`
- Duration: `28800` (8 hours)
- Output: `workday.mp4`
- Result: 960 frames → 32-second video at 30fps

**24-hour day (1-minute intervals)**:
- Interval: `60`
- Duration: `86400` (24 hours)
- Output: `fullday.mp4`
- Result: 1440 frames → 48-second video at 30fps

### Direct Python Usage (Advanced)

You can also run the Python script directly:

```bash
python3 timelapse.py -i 5 -d 300 -o output.mp4

# Arguments:
#   -i, --interval  : Seconds between frames
#   -d, --duration  : Total duration in seconds
#   -o, --output    : Output video path
#   -f, --fps       : Output video FPS (default: 30)
#   -c, --camera    : Camera index (default: 0)
```

## How It Works

### Data Flow

```
1. User Input (Go TUI)
   │
   ├─ Interval: 5 seconds
   ├─ Duration: 600 seconds (10 min)
   └─ Output: timelapse.mp4
   │
   ▼
2. Go spawns subprocess
   │
   └─► python3 timelapse.py -i 5 -d 600 -o timelapse.mp4
       │
       ▼
3. Python captures frames
   │
   ├─ Frame 1 → frames/frame_000001.jpg
   ├─ Frame 2 → frames/frame_000002.jpg
   ├─ ...
   └─ Frame 120 → frames/frame_000120.jpg
   │
   ├─ Progress: [PROGRESS] 1/120 (0.8%)
   ├─ Progress: [PROGRESS] 2/120 (1.7%)
   └─ ...
   │
   ▼
4. Go parses progress
   │
   └─► Updates UI in real-time
   │
   ▼
5. Python compiles video
   │
   └─► ffmpeg -r 30 -i 'frames/*.jpg' -c:v libx264 timelapse.mp4
   │
   ▼
6. Python cleans up
   │
   └─► rm -rf frames/
   │
   ▼
7. Go displays completion
   │
   └─► "Timelapse saved to: timelapse.mp4 ✓"
```

### Communication Protocol

**Python → Go (stdout)**:

```
[INFO] Camera initialized successfully
[PROGRESS] 5/120 (4.2%)
[PROGRESS] 6/120 (5.0%)
[SUCCESS] Video saved to: /path/to/output.mp4
```

**Go parsing**:
- `[PROGRESS]` lines → Extract frame count and percentage
- `[INFO]` / `[ERROR]` → Display in activity log
- `[SUCCESS]` → Show completion message

**Signal handling**:
- User presses `q` → Go sends `SIGINT` to Python process
- Python catches signal → Finishes current frame → Starts compilation
- Python always releases camera and cleans up (in `finally` block)

## Safety Recommendations

### Long-Run Safety

#### 1. Webcam Overheating Prevention

**Problem**: Continuous webcam use can generate heat.

**Solutions**:
- **Limit duration**: Don't exceed 8-12 hours continuous recording
- **Increase interval**: Use 30+ second intervals for long recordings
- **Monitor temperature**: Check laptop temperature periodically
- **Ventilation**: Ensure laptop has good airflow
- **Breaks**: For 24+ hour recordings, split into multiple sessions

```bash
# Example: 24-hour recording split into 3 sessions
# Session 1: Hours 0-8
./chronapse  # Interval: 60, Duration: 28800

# Session 2: Hours 8-16
./chronapse  # Interval: 60, Duration: 28800

# Session 3: Hours 16-24
./chronapse  # Interval: 60, Duration: 28800

# Combine with ffmpeg
ffmpeg -i "concat:part1.mp4|part2.mp4|part3.mp4" -c copy full24h.mp4
```

#### 2. Storage Management

**Calculate required space**:

```
Frames = Duration / Interval
Frame size ≈ 500KB (1080p JPEG)
Total space = Frames × 500KB

Example (8-hour, 30s intervals):
  Frames = 28800 / 30 = 960
  Space = 960 × 500KB ≈ 480MB
```

**Recommendations**:
- Reserve 2× calculated space (safety margin)
- Use `df -h` to check available disk space
- Clean old timelapses regularly
- Monitor disk usage for long recordings:

```bash
watch -n 60 'du -sh frames/ 2>/dev/null || echo "No frames yet"'
```

#### 3. Process Management

**Background execution** (for very long recordings):

```bash
# Using tmux (recommended)
tmux new -s timelapse
./chronapse
# Detach: Ctrl+B, then D
# Reattach: tmux attach -t timelapse

# Using screen
screen -S timelapse
./chronapse
# Detach: Ctrl+A, then D
# Reattach: screen -r timelapse

# Using nohup (basic)
nohup ./chronapse &
```

**Monitor running processes**:

```bash
# Check Python process
ps aux | grep timelapse.py

# Check disk I/O
iotop -p <pid>

# Check camera usage
lsof /dev/video0
```

#### 4. Camera Permission Issues

If you get permission errors:

```bash
# Add user to video group
sudo usermod -a -G video $USER

# Logout and login, or:
newgrp video

# Verify permissions
ls -l /dev/video0
```

#### 5. Error Recovery

**If recording crashes**:

```bash
# Frames are preserved in frames/ directory
# Manually compile:
cd frames
ffmpeg -r 30 -i frame_%06d.jpg -c:v libx264 -pix_fmt yuv420p ../recovered.mp4

# Clean up
cd ..
rm -rf frames/
```

**If camera is stuck**:

```bash
# Kill any hanging processes
pkill -9 -f timelapse.py

# Release camera
sudo rmmod uvcvideo
sudo modprobe uvcvideo
```

#### 6. Power Management

For long recordings:

```bash
# Disable screen sleep
gsettings set org.gnome.desktop.session idle-delay 0

# Disable suspend
sudo systemctl mask sleep.target suspend.target

# Re-enable after recording
gsettings set org.gnome.desktop.session idle-delay 300
sudo systemctl unmask sleep.target suspend.target
```

#### 7. Quality vs. Size Tradeoffs

Adjust FFmpeg settings in `timelapse.py:166-174`:

```python
# Higher quality (larger file)
'-crf', '18',  # Default is 23 (lower = better)

# Lower quality (smaller file)
'-crf', '28',

# Faster encoding (lower quality)
'-preset', 'fast',  # Default is 'medium'

# Better compression (slower)
'-preset', 'slow',
```

## Troubleshooting

### Common Issues

**"Failed to open camera /dev/video0"**:
- Check camera connection: `ls -l /dev/video*`
- Check permissions: `groups` should include `video`
- Close other apps using camera (Chrome, Zoom, etc.)

**"FFmpeg not found"**:
- Install: `sudo apt install ffmpeg`
- Verify: `which ffmpeg`

**"Import cv2 could not be resolved"**:
- Install OpenCV: `pip3 install opencv-python`
- Check: `python3 -c "import cv2"`

**TUI shows garbled characters**:
- Ensure terminal supports UTF-8
- Try different terminal (gnome-terminal, alacritty, etc.)

**Video playback issues**:
- Try VLC: `sudo apt install vlc`
- Reduce FPS in timelapse.py (line 15): `fps=24`

### Debug Mode

Run Python directly to see all output:

```bash
python3 timelapse.py -i 2 -d 20 -o test.mp4
```

Check logs for errors in the Go TUI activity section.

## Performance Tips

1. **Optimal intervals**:
   - Clouds moving: 5-10 seconds
   - Sunset/sunrise: 30-60 seconds
   - Construction: 5-15 minutes
   - Plant growth: 1-5 hours

2. **Output FPS**:
   - 24 FPS: Cinematic look
   - 30 FPS: Smooth motion (default)
   - 60 FPS: Very smooth (larger files)

3. **Resolution** (edit timelapse.py:58-59):
   ```python
   self.camera.set(cv2.CAP_PROP_FRAME_WIDTH, 1920)   # 1080p
   self.camera.set(cv2.CAP_PROP_FRAME_HEIGHT, 1080)

   # 720p for smaller files:
   # self.camera.set(cv2.CAP_PROP_FRAME_WIDTH, 1280)
   # self.camera.set(cv2.CAP_PROP_FRAME_HEIGHT, 720)
   ```

## License

See LICENSE file for details.

## Contributing

Contributions welcome! Please open issues or pull requests on GitHub.
