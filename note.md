
### Các điểm điều chỉnh chính:

#### 1. Thiết kế thanh Top Bar
- **Khung nền Titlebar**: Đồng nhất với background tổng thể của cửa sổ (`linear-gradient(90deg, #1C2426, #202628, #24272A)`), chiều cao chuẩn 38px, loại bỏ border-bottom thô để các khung bên dưới tạo ranh giới tự nhiên.
- **Mac Traffic Lights**: 3 nút macOS (`#FF5F56`, `#FFBD2E`, `#27C93F`) căn lề trái chuẩn xác.
- **Avatar Profile**: Hình chữ nhật bo góc cyan `VT` (`#0891B2`, chữ trắng đậm, padding 1px 5px) cạnh dropdown `VThanh ⌵` và `Version Control ⌵`.
- **Cụm Action trung tâm**: Biểu tượng Database cylinder, Run tròn, Folder dự án và menu `···`.
- **Cụm Utility góc phải**: AI Assistant xoáy, Search kính lúp và Settings bánh răng có **chấm thông báo vàng/amber** ở góc trên (`#E5A122`).

#### 2. Ranh giới tách khung ("Island Cards" layout)
- **Cấu trúc Island Cards**: Cả thanh **Sidebar trái** và **Work Area phải** được đóng gói thành các khối nổi (islands) độc lập:
  - `border-radius: 8px` bo tròn các góc.
  - Viền mỏng tinh tế: `border: 1px solid rgba(255, 255, 255, 0.08)`.
  - Nền khối: `--bg-canvas: #1E1F22`.
- **Khoảng hở phân vùng (`gap: 4px`, `padding: 0 4px 4px 4px`)**: Tạo ranh giới mềm mại giữa Top Bar và các panel, cũng như giữa Database Explorer và Editor tab, đúng phong cách DataGrip New UI.

#### 3. Thanh Tab Bar & Chi tiết Bảng Dữ Liệu
- **Tab Bar**: Bo góc trên 8px ôm vừa vặn khối Work Area; tab active `connectio...credential` có nền xanh đậm `#1F2E4A`, viền xanh sáng `1px solid #3574F0` và nút `×`. Góc phải có nút dropdown `⌵` và menu ba chấm `⋮`.
- **Bảng dữ liệu `connection_credential`**:
  - Thanh lọc `⌵ WHERE` và `≡ ORDER BY` trực tiếp.
  - Các cột: `🔑 id ▽ ⇅`, `📅 createdDate ▽ ⇅`, `📅 lastUpdatedDate ▽ ⇅`, `👤 createdBy ▽`.
  - Giá trị `null` hiển thị dạng `<null>` in nghiêng màu xám mờ.
  - Floating badge góc dưới: `[ 58 rows ⌵ | ⋮ ]`.
  - Terminal console phía dưới hiển thị log thực thi SQL nhiều dòng có highlight cú pháp và thời gian execution/fetching.
- **Status Bar**: Cập nhật breadcrumb điều hướng `Database > [Dev][ReadOnly] 10.220.6.4 > dev_mh_asset > public > tables > connection_credential`, chuông thông báo Gateway và huy hiệu `[Non-commercial use]` xanh lá.