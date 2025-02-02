# Teletrix 🎙️

[![Go Report Card](https://goreportcard.com/badge/github.com/bannawandoor27/Teletrix)](https://goreportcard.com/report/github.com/bannawandoor27/Teletrix)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/bannawandoor27/Teletrix)](https://golang.org/)

Teletrix is a sophisticated real-time audio processing system built in Go, leveraging GStreamer for high-performance audio pipeline management. It provides seamless audio mixing capabilities with voice activity detection.

## 🌟 Features

- **Real-time Audio Processing**: Process audio streams with minimal latency
- **Voice Activity Detection**: Intelligent detection of voice activity in audio streams
- **Audio Mixing**: Mix multiple audio sources (microphone input and test tones)
- **GStreamer Integration**: Utilizes GStreamer's powerful multimedia framework
- **Graceful Shutdown**: Proper resource cleanup and shutdown handling

## 🏗️ Architecture

Teletrix is built with a modular architecture consisting of two main components:

### Audio Processor
- Handles voice activity detection
- Processes audio streams in real-time
- Manages audio state and thresholds

### GStreamer Pipeline
```
[Tone Source (440Hz)] ─┐
                       ├─► [Audio Mixer] ─► [Audio Converter] ─► [Audio Sink]
[Microphone Input] ─┬─► |
                    └─► [Level Monitor]
```

## 🚀 Getting Started

### Prerequisites

- Go 1.16 or higher
- GStreamer and its development files
- Working audio input/output devices

### Installation

1. Install GStreamer and its dependencies:

   ```bash
   # Ubuntu/Debian
   sudo apt-get install libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev

   # macOS
   brew install gstreamer gst-plugins-base gst-plugins-good
   ```

2. Clone the repository:

   ```bash
   git clone https://github.com/teletrix/teletrix.git
   cd teletrix
   ```

3. Install Go dependencies:

   ```bash
   go mod download
   ```

### Usage

1. Build and run the application:

   ```bash
   go build
   ./teletrix
   ```

2. The application will:
   - Initialize the audio processor
   - Set up the GStreamer pipeline
   - Start processing audio input
   - Mix the audio with a 440Hz test tone

3. To stop the application, press `Ctrl+C` for graceful shutdown

## 🛠️ Configuration

The application can be configured through various parameters:

- Audio threshold for voice detection
- Test tone frequency and volume
- Microphone input volume
- Processing intervals

## 🔧 Development

### Project Structure

```
├── internal/
│   ├── audio/
│   │   └── processor.go    # Audio processing and voice detection
│   └── pipeline/
│       └── pipeline.go     # GStreamer pipeline management
└── main.go                 # Application entry point
```

### Key Components

1. **Audio Processor** (`internal/audio/processor.go`)
   - Handles voice activity detection
   - Manages audio processing state
   - Provides clean interface for audio operations

2. **GStreamer Pipeline** (`internal/pipeline/pipeline.go`)
   - Configures and manages the GStreamer pipeline
   - Handles audio mixing and output
   - Provides volume control and pipeline state management

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📚 Documentation

For more detailed documentation about the internal components and their interactions, please refer to the code comments and the following resources:

- [GStreamer Documentation](https://gstreamer.freedesktop.org/documentation/)
- [Go GStreamer Bindings](https://github.com/tinyzimmer/go-gst)

## ✨ Acknowledgments

- GStreamer team for their excellent multimedia framework
- The Go team for the amazing programming language
- Contributors to the go-gst bindings