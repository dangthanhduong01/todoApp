# 🔧 Sửa lỗi: Theme Project không được lưu lại

## 🐛 **Vấn đề trước đây:**
- Project theme (color & background image) không được lưu lại khi chạy lại ứng dụng
- Mỗi lần restart app, theme bị reset về mặc định

## 🔍 **Nguyên nhân:**

### 1. **Lỗi trong `createProject()`**
- Thông tin `BackgroundImage` không được ghi vào file project khi tạo mới
- Chỉ ghi `Color` mà thiếu `BackgroundImage`

### 2. **Lỗi trong `SaveToFile()`** 
- Hàm `SaveToFile()` sử dụng `os.Create()` ghi đè toàn bộ file
- **Mất hết header metadata** (# Project, # Color, # BackgroundImage)
- Chỉ còn lại todo data

### 3. **Lỗi format file không nhất quán**
- File cũ: chỉ có todo data 
- File mới: có header metadata
- Logic load không xử lý được 2 format

## ✅ **Đã sửa:**

### 1. **Fixed `createProject()`** 
```go
// Trước: chỉ ghi Color
content := fmt.Sprintf("# Project: %s\n# Color: %s\n# Created: %s\n\n", name, color, time)

// Sau: ghi cả BackgroundImage
if len(backgroundImage) > 0 && backgroundImage[0] != "" {
    content += fmt.Sprintf("# BackgroundImage: %s\n", backgroundImage[0])
}
```

### 2. **Fixed `SaveToFile()` - Preserve Headers**
```go
// Đọc header trước khi ghi
var headerLines []string
// ... đọc các dòng bắt đầu với #
// Ghi lại header + todo data
```

### 3. **Fixed `LoadFromFile()` - Skip Headers**
```go
// Skip header lines khi load todos
if line == "" || strings.HasPrefix(line, "#") {
    continue // Skip metadata
}
```

### 4. **Added Migration cho file cũ**
```go
// Auto migrate file cũ khi load
func (app *TodoApp) migrateOldProjectFile(projectName string)
```

## 🎯 **Kết quả:**

### ✅ **Theme được lưu và load chính xác**
```
📁 Loaded project: test_theme_project (red) - Background: data/themes/images/sample_background.png
🎨 Applied project theme: test_theme_project (color: red, image: ...)
```

### ✅ **File format đúng**
```
# Project: test_theme_project  
# Color: red
# BackgroundImage: data/themes/images/sample_background.png
# Created: 2025-10-07 15:33:31

1|Test todo|false|2025-10-07T15:37:20+07:00
```

### ✅ **Header được preserve khi thêm/sửa/xóa todo**
- Metadata không bị mất
- Todo data được update bình thường
- Format file nhất quán

## 🔮 **Tương lai:**
- ✅ Tự động migrate file cũ khi user load 
- ✅ Theme persist qua restart
- ✅ Background image hoạt động ổn định
- ✅ No data loss khi thao tác todo

## 📝 **Test Results:**
- ✅ Tạo project với theme → Lưu đúng
- ✅ Restart app → Theme vẫn load
- ✅ Thêm todo → Header không mất  
- ✅ Migration file cũ → Tự động
- ✅ Multiple projects → Theme riêng biệt
