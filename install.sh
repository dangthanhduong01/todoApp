#!/bin/bash

# Todo List App Installation Script
echo "🚀 Installing Todo List App..."

# Get current directory
CURRENT_DIR=$(pwd)
APP_NAME="todoapp"
BINARY_NAME="${APP_NAME}-linux"
DESKTOP_FILE="${APP_NAME}.desktop"
ICON_FILE="${APP_NAME}.png"

# Create applications directory if it doesn't exist
mkdir -p ~/.local/share/applications

# Install directory (you can change this)
INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

echo "📦 Copying binary to $INSTALL_DIR..."
cp "$BINARY_NAME" "$INSTALL_DIR/$APP_NAME"
chmod +x "$INSTALL_DIR/$APP_NAME"

echo "🎨 Installing icon..."
ICON_DIR="$HOME/.local/share/icons"
mkdir -p "$ICON_DIR"
cp "$ICON_FILE" "$ICON_DIR/$ICON_FILE"

echo "🖥️ Creating desktop entry..."
# Update desktop file with correct paths
sed "s|Exec=.*|Exec=$INSTALL_DIR/$APP_NAME|g" "$DESKTOP_FILE" | \
sed "s|Icon=.*|Icon=$ICON_DIR/$ICON_FILE|g" > ~/.local/share/applications/"$DESKTOP_FILE"

chmod +x ~/.local/share/applications/"$DESKTOP_FILE"

echo "🔄 Updating desktop database..."
# Try to update desktop database, fallback if command not available
if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database ~/.local/share/applications/ 2>/dev/null
    echo "   ✅ Desktop database updated"
else
    echo "   ℹ️  update-desktop-database not found (app may need logout/restart to appear in menu)"
fi

echo "✅ Installation completed!"
echo ""
echo "📋 Todo List App has been installed successfully!"
echo "   - Binary location: $INSTALL_DIR/$APP_NAME"
echo "   - Desktop file: ~/.local/share/applications/$DESKTOP_FILE"
echo "   - Icon: $ICON_DIR/$ICON_FILE"
echo ""
echo "🎯 You can now:"
echo "   1. Run from terminal: $APP_NAME (if ~/.local/bin is in PATH)"
echo "   2. Find 'Todo List App' in your application menu"
echo "   3. Run directly: $INSTALL_DIR/$APP_NAME"

# Check if ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo ""
    echo "⚠️  Note: ~/.local/bin is not in your PATH"
    echo "   Add this line to your ~/.bashrc or ~/.zshrc:"
    echo "   export PATH=\"\$HOME/.local/bin:\$PATH\""
fi
