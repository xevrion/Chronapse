# Contributing to Chronapse

Thank you for your interest in contributing to Chronapse! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Enhancements](#suggesting-enhancements)
- [Development Guidelines](#development-guidelines)

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment. We expect all contributors to:

- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Accept constructive criticism gracefully
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

Before you begin contributing:

1. Check existing [issues](https://github.com/xevrion/chronapse/issues) to see if your concern has already been raised
2. For major changes, open an issue first to discuss what you would like to change
3. Fork the repository and create a new branch for your work
4. Keep your changes focused and atomic

## Development Setup

### Prerequisites

Ensure you have the following installed:

```bash
# System dependencies
sudo apt install python3 python3-pip ffmpeg golang-go

# Python dependencies
pip3 install -r requirements.txt

# Go dependencies
go mod download
```

### Building the Project

```bash
# Build the executable
go build -o chronapse main.go

# Run the application
./chronapse
```

## How to Contribute

### Types of Contributions

We welcome various types of contributions:

- **Bug fixes**: Help us squash bugs and improve stability
- **New features**: Implement features from the [roadmap](README.md#todo) or propose new ones
- **Documentation**: Improve or expand documentation
- **Testing**: Add tests to increase code coverage
- **Performance**: Optimize existing code
- **UI/UX**: Enhance the terminal user interface

### Finding Work

- Check the [issue tracker](https://github.com/xevrion/chronapse/issues) for open issues
- Look for issues labeled `good first issue` for beginner-friendly tasks
- Issues labeled `help wanted` are particularly in need of contributors
- Review the [TODO list](README.md#todo) for planned features

## Pull Request Process

### Before Submitting

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR-USERNAME/chronapse.git
   cd chronapse
   ```

2. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

3. **Make Your Changes**
   - Write clean, readable code
   - Follow the coding standards (see below)
   - Test your changes thoroughly
   - Update documentation if necessary

4. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "Brief description of your changes"
   ```

   Follow commit message conventions:
   - Use present tense ("Add feature" not "Added feature")
   - Use imperative mood ("Move cursor to..." not "Moves cursor to...")
   - Limit the first line to 72 characters
   - Reference issues and pull requests when relevant

5. **Push to Your Fork**
   ```bash
   git push origin feature/your-feature-name
   ```

### Submitting the Pull Request

1. Navigate to the original repository
2. Click "New Pull Request"
3. Select your fork and branch
4. Fill out the PR template with:
   - Clear description of changes
   - Related issue numbers (use "Fixes #123" or "Closes #123")
   - Screenshots/recordings for UI changes
   - Testing performed
   - Any breaking changes

### PR Review Process

- Maintainers will review your PR as soon as possible
- Address any requested changes promptly
- Keep the PR focused and avoid scope creep
- Be patient and respectful during the review process
- Once approved, a maintainer will merge your PR

## Coding Standards

### Go Code

- Follow standard Go conventions and formatting
- Use `gofmt` to format your code
- Run `go vet` to catch common issues
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused

Example:
```go
// captureFrame captures a single frame from the webcam
func captureFrame(device string, outputPath string) error {
    // Implementation
}
```

### Python Code

- Follow PEP 8 style guide
- Use type hints where appropriate
- Format code with `black` (if installed)
- Use descriptive variable names
- Add docstrings to functions

Example:
```python
def capture_frame(device: str, output_path: str) -> bool:
    """
    Capture a single frame from the specified device.

    Args:
        device: Path to the camera device
        output_path: Path where the frame should be saved

    Returns:
        True if capture succeeded, False otherwise
    """
    # Implementation
```

### General Guidelines

- Write self-documenting code
- Avoid magic numbers; use named constants
- Handle errors gracefully
- Don't commit commented-out code
- Remove debug print statements
- Keep lines under 100 characters when possible

## Testing

### Manual Testing

Before submitting a PR, test the following scenarios:

1. **Basic functionality**
   - Recording with default settings
   - Recording with custom intervals and durations
   - Early termination with 'q'
   - Invalid input handling

2. **Edge cases**
   - Very short durations (< 10 seconds)
   - Very long intervals
   - Non-existent camera device
   - Disk space issues
   - Permission issues

3. **UI/UX**
   - Keyboard navigation works correctly
   - Progress updates display properly
   - Error messages are clear and helpful
   - Screen resizing doesn't break UI

### Automated Testing

We encourage adding automated tests for new features:

```bash
# For Go code
go test ./...

# For Python code
pytest tests/
```

## Reporting Bugs

When reporting bugs, please use the bug report template and include:

- **Description**: Clear and concise description of the bug
- **Steps to Reproduce**: Numbered steps to reproduce the behavior
- **Expected Behavior**: What you expected to happen
- **Actual Behavior**: What actually happened
- **Environment**:
  - OS and version (e.g., Ubuntu 22.04)
  - Go version (`go version`)
  - Python version (`python3 --version`)
  - FFmpeg version (`ffmpeg -version`)
  - OpenCV version
- **Screenshots/Logs**: If applicable, add screenshots or error logs
- **Additional Context**: Any other relevant information

## Suggesting Enhancements

When suggesting enhancements, please use the feature request template and include:

- **Problem Statement**: What problem does this solve?
- **Proposed Solution**: Detailed description of your proposed feature
- **Alternatives Considered**: Other solutions you've considered
- **Use Cases**: Real-world scenarios where this would be useful
- **Implementation Ideas**: Technical approach (if you have one)
- **Mockups**: UI mockups or diagrams (if applicable)

## Development Guidelines

### Working with the TUI

- The TUI uses [Bubbletea](https://github.com/charmbracelet/bubbletea)
- Follow the Elm architecture pattern (Model-Update-View)
- Keep state changes predictable and testable
- Use Lipgloss for styling

### Working with the Backend

- The Python backend communicates via stdout/stderr
- Progress updates use the format: `[PROGRESS] current/total (percent%)`
- Log messages are plain text
- Exit codes indicate success/failure

### Cross-Component Changes

If your changes affect both Go and Python components:
- Test the integration thoroughly
- Document the communication protocol changes
- Ensure backward compatibility when possible

### Performance Considerations

- Minimize disk I/O operations
- Use efficient algorithms for frame processing
- Consider memory usage for long recordings
- Profile code for bottlenecks

### Security Considerations

- Validate all user inputs
- Avoid command injection vulnerabilities
- Handle file paths safely
- Don't expose sensitive information in logs

## Questions?

If you have questions about contributing:

1. Check existing documentation (README.md, USAGE.md)
2. Search closed issues for similar questions
3. Open a new issue with the `question` label
4. Be specific and provide context

## Recognition

Contributors will be recognized in the project. Significant contributions may be highlighted in release notes.

Thank you for contributing to Chronapse!
