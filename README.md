# 📝 Todo List Desktop Application

Ứng dụng quản lý công việc desktop được viết bằng Go và sử dụng framework Fyne cho giao diện người dùng.

## ✨ Tính năng

- **Thêm công việc mới**: Nhập mô tả công việc và thêm vào danh sách
- **Giao diện Tab**: 3 tab riêng biệt - Tất cả, Chưa hoàn thành, Đã hoàn thành
- **Nút hành động riêng**: Mỗi công việc có nút ✅ Hoàn thành và 🗑️ Xóa riêng
- **Xem danh sách công việc**: Hiển thị công việc theo trạng thái với emoji rõ ràng
- **Đánh dấu hoàn thành**: Click nút ✅ bên cạnh mỗi công việc
- **Xóa công việc**: Click nút 🗑️ với xác nhận trước khi xóa
- **Lưu trữ bền vững**: Dữ liệu được lưu trong file text (`todos.txt`)
- **Giao diện card**: Mỗi công việc hiển thị dạng card với thông tin rõ ràng

## 🚀 Cách sử dụng

### Yêu cầu hệ thống
- Go 1.19 hoặc cao hơn
- Linux với X11 (hoặc Wayland với XWayland)
- Các thư viện hệ thống: libgl1-mesa-dev, libxi-dev, libxcursor-dev, libxrandr-dev, libxinerama-dev, libxxf86vm-dev

### Cài đặt dependencies
```bash
sudo apt update
sudo apt install -y libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev libxxf86vm-dev
```

### Build ứng dụng
```bash
go build
```

### Chạy ứng dụng
```bash
./todoapp
```

## 📁 Cấu trúc dự án

```
todoapp/
├── main.go          # Giao diện người dùng với Fyne
├── todo.go          # Logic quản lý todos và file operations
├── todos.txt        # File lưu trữ dữ liệu (tự động tạo)
├── go.mod           # Go module dependencies
└── README.md        # Tài liệu này
```

## 🛠️ Phát triển

### Dependencies chính
- `fyne.io/fyne/v2` - Framework GUI cho Go
- Go standard library cho file I/O và string processing

### Định dạng dữ liệu
Dữ liệu được lưu trong file `todos.txt` với format:
```
ID|Description|Completed|CreatedAt
1|Mua sữa|false|2024-01-01T10:00:00Z
2|Làm bài tập|true|2024-01-01T11:00:00Z
```

## 🔧 Tùy chỉnh

Bạn có thể tùy chỉnh:
- Đường dẫn file lưu trữ trong hàm `main()`
- Giao diện người dùng trong các hàm `setupUI()`
- Thêm các tính năng mới như filter, search, priority, v.v.

## 📱 Giao diện

Ứng dụng có giao diện hiện đại với:
- **Header**: Tiêu đề ứng dụng và mô tả
- **Input Section**: Field nhập và nút thêm công việc mới
- **Tab Navigation**: 3 tab để lọc theo trạng thái
  - 📋 **Tất cả**: Hiển thị toàn bộ công việc
  - 📌 **Chưa hoàn thành**: Chỉ công việc đang thực hiện  
  - ✅ **Đã hoàn thành**: Chỉ công việc đã xong
- **Todo Cards**: Mỗi công việc hiển thị dạng card với:
  - Emoji trạng thái (📌/✅) và mô tả công việc
  - ID công việc để dễ theo dõi
  - Nút ✅ Hoàn thành (hoặc thông báo nếu đã hoàn thành)
  - Nút 🗑️ Xóa với xác nhận
- **Dialogs**: Thông báo xác nhận và trạng thái
