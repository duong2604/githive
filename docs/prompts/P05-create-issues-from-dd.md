---
id: P05
title: Create GitHub issues from approved DDs
step: "3"
status: Ready
branch: docs/issues-from-dd
pr: ""
depends_on: [P04]
decisions: []
last_updated: 2026-11-13
---

# P05 — Tạo issue từ DD đã duyệt

**Điều kiện chạy:** PR của P02 đã merge, và DD-00, DD-03 đã ở trạng thái `Approved`.

## Prompt (dán nguyên văn)

```text
Prompt này là docs/prompts/P05-create-issues-from-dd.md. Đọc CLAUDE.md và docs/PROGRESS.md. Với mỗi file trong docs/detail-design/ có status: Approved mà chưa có issue (xem bảng Traceability trong docs/requirements/README.md):
- Tạo milestone MVP và các label (mvp, v1, v2, design, bug, windows, wsl) nếu chưa có.
- Tạo issue tiêu đề "[DD-xx] <tên>". Body gồm: link tới DD, Acceptance criteria chép đúng từ mục acceptance/test của DD (không tự thêm bớt), danh sách package bị ảnh hưởng, và mục "Can run in parallel with".
- DD-00 đã tách sẵn thành DD-00a, DD-00b, DD-00c (DD-00 §1): tạo một issue cho mỗi phần, không hỏi lại. DD khác quá lớn (hơn khoảng 400 dòng code dự kiến) thì đề xuất cách tách thành nhiều issue và hỏi tôi trước khi tạo.
- Mở một PR docs trên nhánh docs/issues-from-dd: cập nhật cột Issue trong bảng Traceability và trong docs/basic-design/README.md, cập nhật docs/PROGRESS.md (bước 3 → In progress) và dòng P05 trong docs/prompts/README.md.
Kết thúc bằng báo cáo theo mẫu CLAUDE.md §13, tiêu đề "Báo cáo — P05 — <ngày>". Mục 3 liệt kê các issue và chỉ rõ issue nào giao được cho nhiều agent làm song song. Không bắt đầu code.
```

## Sau P05

Mỗi issue được giao bằng một prompt riêng, theo các phần trong DD-00 §1:

1. **DD-00a** làm trước: `P06` (DD-00a).
2. Sau khi DD-00a merge, ba issue chạy song song: `P07` (**DD-00b**), `P08` (**DD-00c**), `P09` (**DD-03**).

Cố vấn soạn các prompt này dựa trên mẫu C trong `P00-kickoff.md`, mở đầu bằng dòng "Implements DD-xx" và kết thúc bằng yêu cầu báo cáo theo §13.
