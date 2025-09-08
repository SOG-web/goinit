#!/bin/bash

# Go Gin API Generator Installation Script

echo "🚀 Installing Go Gin API Generator..."
echo "====================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

# Build the generator
echo "📦 Building generator..."
cd cli-generator
go build -o goinit-generator main.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

# Create installation directory
INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

# Copy binary
cp goinit-generator "$INSTALL_DIR/"

echo "✅ Installation complete!"
echo ""
echo "📖 Usage:"
echo "   goinit-generator"
echo ""
echo "Or add to PATH:"
echo "   export PATH=\"$INSTALL_DIR:\$PATH\""
echo ""
echo "🎉 Happy coding!"