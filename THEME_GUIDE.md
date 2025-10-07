# Hướng dẫn sử dụng Theme cho Projects

## Tính năng mới: Project Theme với Background Image

### Cách sử dụng:

1. **Tạo Project mới với Theme:**
   - Vào tab "📁 Projects"
   - Nhấn nút "+ Tạo Project"
   - Nhập tên project
   - Chọn màu chủ đạo
   - Nhấn "📁 Chọn ảnh nền" để chọn ảnh background
   - Nhấn "Tạo" để tạo project

2. **Định dạng ảnh hỗ trợ:**
   - `.jpg`, `.jpeg`
   - `.png`
   - `.gif`
   - `.bmp`

3. **Cách ảnh được lưu trữ:**
   - Ảnh được copy vào thư mục `data/themes/images/`
   - Thông tin theme được lưu trong file project (`.txt`)
   - Mỗi project có thể có theme riêng biệt

4. **Tab Todos:**
   - Vẫn giữ nguyên giao diện và theme như cũ
   - Không bị ảnh hưởng bởi theme của project

5. **Tab Projects:**
   - Hiển thị background image của project đang active
   - Màu chủ đề áp dụng theo từng project
   - Theme thay đổi khi chuyển project

### Cấu trúc file project:
```
# Project: Tên project
# Color: blue
# BackgroundImage: data/themes/images/background.jpg
# Created: 2025-01-01 12:00:00

Todo items...
```

### Lưu ý:
- Theme chỉ áp dụng cho tab Projects
- Tab Todos vẫn sử dụng theme chung của ứng dụng
- Mỗi project có thể có màu và ảnh nền riêng
- Ảnh sẽ được copy vào thư mục project để đảm bảo tính ổn định
