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
rm -f "$HOME/.local/share/icons/$ICON_FILE"
rm -f "$HOME/.local/share/icons/hicolor/128x128/apps/$ICON_FILE"

echo "🖥️ Removing desktop entry..."
rm -f ~/.local/share/applications/"$DESKTOP_FILE"

echo "🖥️ Removing desktop shortcut..."
rm -f "$HOME/Desktop/$DESKTOP_FILE"

echo "🔄 Updating desktop database..."
if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database ~/.local/share/applications/ 2>/dev/null
    echo "   ✅ Desktop database updated"
else
    echo "   ℹ️  update-desktop-database not found"
fi

# Also refresh icon cache
if command -v gtk-update-icon-cache >/dev/null 2>&1; then
    gtk-update-icon-cache ~/.local/share/icons/ 2>/dev/null || true
fi

echo "✅ Uninstallation completed!"
echo "📋 Todo List App has been removed from your system."
