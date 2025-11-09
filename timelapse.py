#!/usr/bin/env python3
"""
Chronapse Timelapse Recorder
Captures frames from webcam at specified intervals and compiles them into a video.
"""

import argparse
import cv2
import os
import sys
import signal
import time
import subprocess
import shutil
from pathlib import Path
from datetime import datetime


class TimelapseRecorder:
    """Handles webcam frame capture and video compilation."""

    def __init__(self, interval, duration, output_path, fps=30, camera_index=0):
        """
        Initialize the timelapse recorder.

        Args:
            interval: Seconds between each frame capture
            duration: Total recording duration in seconds
            output_path: Path for the output video file
            fps: Frames per second for the output video
            camera_index: Webcam device index (default 0 for /dev/video0)
        """
        self.interval = interval
        self.duration = duration
        self.output_path = Path(output_path)
        self.fps = fps
        self.camera_index = camera_index
        self.frames_dir = Path("frames")
        self.camera = None
        self.should_stop = False
        self.frames_captured = 0

        # Calculate expected frame count
        self.expected_frames = int(duration / interval)

        # Setup signal handler for graceful shutdown
        signal.signal(signal.SIGINT, self._signal_handler)
        signal.signal(signal.SIGTERM, self._signal_handler)

    def _signal_handler(self, signum, frame):
        """Handle interrupt signals gracefully."""
        print("\n[INFO] Received stop signal. Finishing recording...", flush=True)
        self.should_stop = True

    def _initialize_camera(self):
        """Initialize the webcam connection."""
        print(f"[INFO] Initializing camera /dev/video{self.camera_index}...", flush=True)

        # Try to open the camera with V4L2 backend on Linux
        self.camera = cv2.VideoCapture(self.camera_index, cv2.CAP_V4L2)

        if not self.camera.isOpened():
            raise RuntimeError(f"Failed to open camera /dev/video{self.camera_index}")

        # Set camera properties for better quality
        self.camera.set(cv2.CAP_PROP_FRAME_WIDTH, 1920)
        self.camera.set(cv2.CAP_PROP_FRAME_HEIGHT, 1080)
        self.camera.set(cv2.CAP_PROP_FPS, 30)

        # Warm up camera with a test frame
        ret, _ = self.camera.read()
        if not ret:
            raise RuntimeError("Failed to capture test frame from camera")

        print("[INFO] Camera initialized successfully", flush=True)

    def _setup_frames_directory(self):
        """Create a clean frames directory."""
        if self.frames_dir.exists():
            print(f"[INFO] Cleaning existing frames directory...", flush=True)
            shutil.rmtree(self.frames_dir)

        self.frames_dir.mkdir(exist_ok=True)
        print(f"[INFO] Created frames directory: {self.frames_dir}", flush=True)

    def _capture_frame(self, frame_number):
        """
        Capture a single frame from the camera.

        Args:
            frame_number: Sequential frame number for filename

        Returns:
            bool: True if capture was successful
        """
        ret, frame = self.camera.read()

        if not ret:
            print(f"[ERROR] Failed to capture frame {frame_number}", flush=True)
            return False

        # Save frame with zero-padded filename for proper sorting
        frame_path = self.frames_dir / f"frame_{frame_number:06d}.jpg"
        cv2.imwrite(str(frame_path), frame, [cv2.IMWRITE_JPEG_QUALITY, 95])

        self.frames_captured += 1

        # Progress output that Go TUI can parse
        progress = (self.frames_captured / self.expected_frames) * 100
        print(f"[PROGRESS] {self.frames_captured}/{self.expected_frames} ({progress:.1f}%)", flush=True)

        return True

    def _record_frames(self):
        """Main recording loop."""
        print(f"[INFO] Starting recording: {self.expected_frames} frames over {self.duration}s", flush=True)
        print(f"[INFO] Capture interval: {self.interval}s", flush=True)

        start_time = time.time()
        frame_number = 0

        while not self.should_stop:
            current_time = time.time()
            elapsed = current_time - start_time

            # Check if we've reached the duration limit
            if elapsed >= self.duration:
                print("[INFO] Reached target duration", flush=True)
                break

            # Capture frame
            if self._capture_frame(frame_number):
                frame_number += 1

            # Calculate next capture time
            next_capture_time = start_time + (frame_number * self.interval)
            sleep_time = next_capture_time - time.time()

            # Sleep until next capture, checking for stop signal
            if sleep_time > 0:
                # Break sleep into smaller chunks to respond faster to signals
                sleep_chunks = int(sleep_time / 0.1) + 1
                for _ in range(sleep_chunks):
                    if self.should_stop:
                        break
                    time.sleep(min(0.1, sleep_time))
                    sleep_time -= 0.1

        actual_duration = time.time() - start_time
        print(f"[INFO] Recording complete: {self.frames_captured} frames in {actual_duration:.1f}s", flush=True)

    def _compile_video(self):
        """Compile captured frames into a video using ffmpeg."""
        if self.frames_captured < 2:
            print("[ERROR] Not enough frames to create video (minimum 2 required)", flush=True)
            return False

        print(f"[INFO] Compiling video with {self.frames_captured} frames at {self.fps} FPS...", flush=True)

        # Ensure output directory exists
        self.output_path.parent.mkdir(parents=True, exist_ok=True)

        # Build ffmpeg command
        ffmpeg_cmd = [
            # 'ffmpeg',
            '/usr/bin/ffmpeg',
            '-y',  # Overwrite output file if it exists
            '-framerate', str(self.fps),
            '-pattern_type', 'glob',
            '-i', str(self.frames_dir / 'frame_*.jpg'),
            '-c:v', 'libx264',
            # '-preset', 'medium',
            '-crf', '23',  # Quality factor (lower = better quality)
            '-pix_fmt', 'yuv420p',  # Compatibility with most players
            str(self.output_path)
        ]

        try:
            # Run ffmpeg with suppressed output except errors
            result = subprocess.run(
                ffmpeg_cmd,
                capture_output=True,
                text=True,
                check=True
            )

            print(f"[SUCCESS] Video saved to: {self.output_path.absolute()}", flush=True)

            # Get output file size
            file_size = self.output_path.stat().st_size / (1024 * 1024)  # MB
            print(f"[INFO] Output file size: {file_size:.2f} MB", flush=True)

            return True

        except subprocess.CalledProcessError as e:
            print(f"[ERROR] FFmpeg failed: {e.stderr}", flush=True)
            return False
        except FileNotFoundError:
            print("[ERROR] FFmpeg not found. Please install it: sudo apt install ffmpeg", flush=True)
            return False

    def _cleanup_frames(self):
        """Remove temporary frames directory."""
        if self.frames_dir.exists():
            print(f"[INFO] Cleaning up frames directory...", flush=True)
            shutil.rmtree(self.frames_dir)
            print("[INFO] Cleanup complete", flush=True)

    def _release_camera(self):
        """Release the camera resource."""
        if self.camera is not None and self.camera.isOpened():
            self.camera.release()
            print("[INFO] Camera released", flush=True)

    def run(self):
        """Execute the complete timelapse recording workflow."""
        try:
            # Setup phase
            self._setup_frames_directory()
            self._initialize_camera()

            # Recording phase
            self._record_frames()

            # Always release camera before compilation
            self._release_camera()

            # Compilation phase
            if self.frames_captured > 0:
                success = self._compile_video()

                # Cleanup phase
                self._cleanup_frames()

                return 0 if success else 1
            else:
                print("[ERROR] No frames captured", flush=True)
                self._cleanup_frames()
                return 1

        except Exception as e:
            print(f"[ERROR] Unexpected error: {e}", flush=True)
            return 1
        finally:
            # Ensure camera is always released
            self._release_camera()


def parse_arguments():
    """Parse command-line arguments."""
    parser = argparse.ArgumentParser(
        description='Chronapse: Webcam Timelapse Recorder',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter
    )

    parser.add_argument(
        '-i', '--interval',
        type=float,
        required=True,
        help='Seconds between each frame capture'
    )

    parser.add_argument(
        '-d', '--duration',
        type=float,
        required=True,
        help='Total recording duration in seconds'
    )

    parser.add_argument(
        '-o', '--output',
        type=str,
        required=True,
        help='Output video file path'
    )

    parser.add_argument(
        '-f', '--fps',
        type=int,
        default=30,
        help='Output video frames per second'
    )

    parser.add_argument(
        '-c', '--camera',
        type=int,
        default=0,
        help='Camera device index'
    )

    return parser.parse_args()


def main():
    """Main entry point."""
    print(f"[INFO] Chronapse Timelapse Recorder started at {datetime.now()}", flush=True)

    args = parse_arguments()

    # Validate arguments
    if args.interval <= 0:
        print("[ERROR] Interval must be positive", flush=True)
        return 1

    if args.duration <= 0:
        print("[ERROR] Duration must be positive", flush=True)
        return 1

    if args.fps <= 0:
        print("[ERROR] FPS must be positive", flush=True)
        return 1

    # Create recorder and run
    recorder = TimelapseRecorder(
        interval=args.interval,
        duration=args.duration,
        output_path=args.output,
        fps=args.fps,
        camera_index=args.camera
    )

    exit_code = recorder.run()

    print(f"[INFO] Chronapse finished with exit code {exit_code}", flush=True)
    return exit_code


if __name__ == '__main__':
    sys.exit(main())
