#!/bin/bash

# Exit on error
set -e

echo "🚀 Starting setup for Recorder project..."

# Check if running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    echo "❌ This setup script is only for macOS"
    exit 1
fi

# Check for Homebrew
if ! command -v brew &> /dev/null; then
    echo "📦 Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

# Check for Go
if ! command -v go &> /dev/null; then
    echo "📦 Installing Go..."
    brew install go
fi

# Check for Python
if ! command -v python3 &> /dev/null; then
    echo "📦 Installing Python..."
    brew install python
fi

# Check for BlackHole-2ch
if ! brew list --cask blackhole-2ch &> /dev/null; then
    echo "📦 Installing BlackHole-2ch audio driver..."
    brew install --cask blackhole-2ch
    echo "⚠️ Please restart your computer after installation to enable BlackHole-2ch"
fi

# Check for CMake
if ! command -v cmake &> /dev/null; then
    echo "📦 Installing CMake..."
    brew install cmake
fi

# Check for ffmpeg
if ! command -v ffmpeg &> /dev/null; then
    echo "📦 Installing ffmpeg..."
    brew install ffmpeg
fi

# Clone and build whisper.cpp if not already present
if [ ! -d "whisper.cpp" ]; then
    echo "📦 Cloning whisper.cpp..."
    git clone https://github.com/ggerganov/whisper.cpp.git
    cd whisper.cpp
    echo "📦 Building whisper.cpp..."
    make
    cd ..
fi

# Download whisper model if not present
if [ ! -f "whisper.cpp/models/ggml-base.bin" ]; then
    echo "📦 Downloading whisper base model..."
    cd whisper.cpp
    ./models/download-ggml-model.sh base
    cd ..
fi

# Install Python dependencies
echo "📦 Installing Python dependencies..."
pip3 install -r requirements.txt

echo "✅ Setup completed successfully!"
echo "📝 Next steps:"
echo "1. Make sure BlackHole-2ch is selected as your system output device in System Preferences > Sound > Output"
echo "2. Run the recorder with: go run cmd/record.go"
echo "3. To transcribe audio, use: python3 scripts/transcribe.py" 