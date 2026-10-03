Listed directory table
Viewed Table.svelte:1-18
Viewed styles.css:1-35
Searched for "<table"
Viewed TableDataView.svelte:230-285
Viewed AuditView.svelte:110-165

Lý do trong trang admin hiện tại chưa sử dụng core component `Table` từ `@ohjanus/ui` đến từ **sự khác biệt cốt lõi về bản chất nghiệp vụ (Domain Architecture)** giữa một **Web Table thông thường** và một **Database IDE Data Grid**:

---

### 1. Bản chất: "General Web Table" vs. "IDE Database Data Grid"

| Tiêu chí | Core Component `Table` (`packages/ui`) | Bảng trong Admin (`TableDataView`, `SqlConsoleView`) |
| :--- | :--- | :--- |
| **Mục đích** | Hiển thị dữ liệu web/danh mục thông thường (User, Settings, Post...) | Mô phỏng **DataGrip / DBeaver / Excel spreadsheet** chuyên nghiệp |
| **Cấu trúc ô** | Border-bottom phân cách dòng, padding thoáng web | **Lưới ô 2 chiều (2D Grid)**: Mỗi cell có viền 4 cạnh (`border-right`, `border-bottom`) tạo thành từng ô tính toán |
| **Row Index** | Không có | Có cột **Row Number Gutter** cố định (`1, 2, 3...`) để theo dõi số thứ tự dòng và chọn dòng |
| **Font & Alignment** | Font UI (`Inter`), padding `6px 10px` | Font Monospace (`JetBrains Mono`), tabular nums, chiều cao dòng cố định tuyệt đối (`--table-row-height: 28px/30px`) |

---

### 2. Các tính năng chuyên biệt của DataGrip mà `Table` cơ bản chưa hỗ trợ

1. **Header tích hợp đa tương tác (Multi-action Column Header)**:
   - Trong `TableDataView` và `SqlConsoleView`, mỗi `<th>` không chỉ chứa text mà còn có:
     - Icon kiểu dữ liệu SQL (Primary Key 🔑, Clock 🕒, User 👤, Hash #).
     - Nút phễu lọc nhanh theo từng cột (`header-funnel`).
     - Nút sắp xếp ASC/DESC (`header-sort`).
     - Kéo giãn kích thước cột (Column Resizing).
2. **Cơ chế Sticky Header 2 chiều**:
   - Khi bảng kết quả truy vấn có hàng nghìn dòng và hàng chục cột, header phải `position: sticky; top: 0` cùng z-index chuẩn, đồng thời cột STT phải cố định bên trái khi scroll ngang.
3. **Đặc thù hiển thị kiểu dữ liệu Database**:
   - Giá trị `NULL` hiển thị in nghiêng mờ (`<null>`).
   - Cột UUID rút gọn, boolean indicator, JSON viewer, foreign key link.
   - Thao tác click chọn từng ô (Cell selection) và chỉnh sửa trực tiếp (Inline Cell Editing).
4. **Component Overhead & Performance**:
   - Core `Table` hiện tại là dạng compound component (`Table` > `TableHeader` > `TableRow` > `TableCell`).
   - Khi query SQL trả về **hàng trăm đến hàng ngàn rows**, việc khởi tạo hàng chục nghìn Svelte component instance con cho từng cell (`<TableCell>`) sẽ gây ra runtime overhead về memory và render tree so với việc render trực tiếp thẻ HTML `<table>`, `<tr>`, `<td>` tối ưu.

---

### 3. Đánh giá & Hướng chuẩn hóa chuẩn Senior

Theo tiêu chuẩn kiến trúc các công cụ database (như DataGrip, Supabase Studio, VS Code Data Table):

1. **Tách biệt 2 loại component trong Design System**:
   - `Table`: Dành cho bảng danh mục quản trị đơn giản (chỉ có dòng và cột cơ bản).
   - `DataGrid` (hoặc `DatabaseTable`): Dành cho bảng dữ liệu database chuyên sâu (hỗ trợ sticky header, row number, inline edit, cell selection).
2. **Những nơi trong Admin NÊN dùng core `Table`**:
   - Trang **[TokensView.svelte](file:///Users/thanhtran/Desktop/myapps/open-sources/ohjanus/ui/src/lib/components/views/TokensView.svelte)** (bảng danh sách API Tokens) là dữ liệu dạng danh mục quản trị đơn giản, rất phù hợp để refactor sang dùng component `Table` từ `@ohjanus/ui`.
3. **Những nơi NÊN giữ dạng Data Grid**:
   - **[TableDataView.svelte](file:///Users/thanhtran/Desktop/myapps/open-sources/ohjanus/ui/src/lib/components/views/TableDataView.svelte)** và **[SqlConsoleView.svelte](file:///Users/thanhtran/Desktop/myapps/open-sources/ohjanus/ui/src/lib/components/views/SqlConsoleView.svelte)**: Cần giữ nguyên kiến trúc Data Grid chuyên biệt (hoặc trừu tượng hóa thành `@ohjanus/ui/components/data-grid`) để bảo toàn đầy đủ trải nghiệm DataGrip.