<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Todo List Desktop Application - Copilot Instructions

This is a Go desktop application using the Fyne framework for GUI development. 

## Project Context
- **Language**: Go (Golang)
- **GUI Framework**: Fyne v2
- **Data Storage**: Plain text file (todos.txt)
- **Target Platform**: Linux desktop (with X11/Wayland support)

## Key Components
1. **main.go**: Contains the GUI application using Fyne framework
2. **todo.go**: Contains the TodoList struct and file I/O operations
3. **todos.txt**: Data persistence file (format: ID|Description|Completed|CreatedAt)

## Development Guidelines
- Use Fyne v2 widgets and containers for UI components
- Follow Go conventions for naming and code organization
- Handle errors gracefully with user-friendly dialog messages
- Keep data persistence simple with text file format
- Use Vietnamese text for user interface elements
- Maintain clean separation between GUI logic and business logic

## Common Patterns
- Use `widget.New*()` functions to create Fyne widgets
- Use `container.New*()` for layout containers
- Use `dialog.Show*()` for user interactions and confirmations
- Handle file I/O errors appropriately
- Use Go time.Time for timestamp management

## Dependencies
- fyne.io/fyne/v2/app
- fyne.io/fyne/v2/widget  
- fyne.io/fyne/v2/container
- fyne.io/fyne/v2/dialog
- fyne.io/fyne/v2/theme
