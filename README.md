# 📝 Todo List Desktop Application

This desktop task management application is written in Go and uses the Fyne framework for the user interface.

## ✨ Features

- **Add New Task**: Enter the task description and add it to the list.
- **Tab Interface**: 3 separate tabs - All, Not Completed, Completed.
- **Individual Action Buttons**: Each task has its own ✅ Complete and 🗑️ Delete buttons.
- **View Task List**: Displays tasks by status with clear emojis.
- **Mark as Complete**: Click the ✅ button next to each task.
- **Delete Task**: Click the 🗑️ button with confirmation before deleting.
- **Persistent Storage**: Data is stored in a text file (`todos.txt`).
- **Card Interface**: Each task is displayed as a card with clear information.

## 🚀 How to use

### System Requirements
- Go 1.19 or higher
- Linux with X11 (or Wayland with XWayland)
- Các thư viện hệ thống: libgl1-mesa-dev, libxi-dev, libxcursor-dev, libxrandr-dev, libxinerama-dev, libxxf86vm-dev

### Setup dependencies
```bash
sudo apt update
sudo apt install -y libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev libxxf86vm-dev
```

### Build
```bash
go build
```

### Run
```bash
./todoapp
```

## 📁 Project structure

```
todoapp/
├── main.go          # Giao diện người dùng với Fyne
├── todo.go          # Logic quản lý todos và file operations
├── todos.txt        # File lưu trữ dữ liệu (tự động tạo)
├── go.mod           # Go module dependencies
└── README.md        # Tài liệu này
```

## 🛠️ Development

### Dependencies
- `fyne.io/fyne/v2` - Framework GUI cho Go
- Go standard library cho file I/O và string processing

### Data type
Data store in file `todos.txt` with format:
```
ID|Description|Completed|CreatedAt
1|Mua sữa|false|2024-01-01T10:00:00Z
2|Làm bài tập|true|2024-01-01T11:00:00Z
```

## 🔧 Customization
You can customize:
- The file path stored in the `main()` function
- The user interface in the `setupUI()` functions
- Add new features such as filters, search, priority, etc.

## 📱 Interface
- **Header**: Application title and description
- **Input Section**: Input field and add new task button
- **Navigation Tabs**: 3 tabs to filter by status
  - 📋 **All**: Displays all tasks
  - 📌 **Incomplete**: Shows tasks currently in progress
  - ✅ **Completed**: Shows completed tasks
- **Todo Cards**: Each task is displayed as a card with:
  - Status emoji (📌/✅) and task description
  - Task ID for easy tracking
  - ✅ Complete button (or notification if completed)
  - 🗑️ Delete button with confirmation
- **Dialogs**: Confirmation and status messages


## Todo:
- Add a Career Pathway roadmap tab.
