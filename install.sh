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
# Install to multiple icon directories for compatibility
ICON_BASE_DIR="$HOME/.local/share/icons"

# Install to hicolor theme (standard)
HICOLOR_DIR="$ICON_BASE_DIR/hicolor/128x128/apps"
mkdir -p "$HICOLOR_DIR"
cp "$ICON_FILE" "$HICOLOR_DIR/$ICON_FILE"

# Also install to base icons directory
mkdir -p "$ICON_BASE_DIR"
cp "$ICON_FILE" "$ICON_BASE_DIR/$ICON_FILE"


echo "🖥️ Creating desktop entry..."
# Create desktop file dynamically with current project path
cat > ~/.local/share/applications/"$DESKTOP_FILE" << EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=Todo List App
Comment=Desktop Todo List Application with Project Management
Exec=bash -c "cd '$CURRENT_DIR' && ./$APP_NAME"
Icon=$CURRENT_DIR/$ICON_FILE
Terminal=false
Categories=Office;Utility;
Keywords=todo;task;project;productivity;
StartupNotify=true
StartupWMClass=Todo List Application
X-Ubuntu-Touch=true
EOF

echo "🔄 Updating desktop database..."
# Update desktop database
if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database ~/.local/share/applications/ 2>/dev/null
    echo "   ✅ Desktop database updated"
fi
