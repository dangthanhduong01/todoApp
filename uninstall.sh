#!/bin/bash

# Todo List App Uninstallation Script
echo "🗑️ Uninstalling Todo List App..."

APP_NAME="todoapp"
DESKTOP_FILE="${APP_NAME}.desktop"
ICON_FILE="${APP_NAME}.png"

INSTALL_DIR="$HOME/.local/bin"
ICON_DIR="$HOME/.local/share/icons"

echo "📦 Removing binary..."
rm -f "$INSTALL_DIR/$APP_NAME"

echo "🎨 Removing icon..."
rm -f "$ICON_DIR/$ICON_FILE"

echo "🖥️ Removing desktop entry..."
rm -f ~/.local/share/applications/"$DESKTOP_FILE"

echo "🔄 Updating desktop database..."
update-desktop-database ~/.local/share/applications/ 2>/dev/null

echo "✅ Uninstallation completed!"
echo "📋 Todo List App has been removed from your system."
