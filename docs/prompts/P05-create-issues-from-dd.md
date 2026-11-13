---
id: P05
title: Create GitHub issues from approved DDs
step: "3"
status: Ready
branch: docs/issues-from-dd
pr: ""
depends_on: [P04, P04b]
decisions: []
last_updated: 2026-11-13
---

# P05 — Tạo issue từ DD đã duyệt

**Điều kiện chạy:** PR #6 (P03, P04, P04b) đã merge vào `main`, nên DD-00 và DD-03 rev 2 đã ở trạng thái `Approved` trên `main`.

## Việc Danny làm trước

**1. Chuẩn bị máy** (chạy ở gốc repo, trên máy nào cũng được):

```bash
git switch main && git pull --ff-only
```

**2. Đặt file:** không có file nào cần chép. P05 chỉ đọc các DD đã có trên `main`.

**3. Kết quả `git status --short` phải rỗng** (không có dòng nào). Nếu còn thay đổi chưa commit, xử lý xong rồi mới chạy P05.

**4.** Mở Claude Code ở gốc repo, dán phần "Prompt (dán nguyên văn)" bên dưới.

## Danh sách issue (P04b-D11)

| # | Issue | DD | Phụ thuộc | Nhóm chạy |
| --- | --- | --- | --- | --- |
| 1 | `[DD-00a] Errors, output, modes and platform helpers` | DD-00 §1 | — | Nhóm 1 (một mình) |
| 2 | `[DD-00b] Git runner and version check` | DD-00 §1 | DD-00a | Nhóm 2 (song song) |
| 3 | `[DD-00c] Repository lock` | DD-00 §1 | DD-00a | Nhóm 2 (song song) |
| 4 | `[DD-03a] Target parser and resolver` | DD-03 §1 | DD-00a | Nhóm 2 (song song) |
| 5 | `[DD-03b] Current-directory locator on git` | DD-03 §1 | DD-00b, DD-03a | Nhóm 3 |

## Prompt (dán nguyên văn)

```text
Prompt này là docs/prompts/P05-create-issues-from-dd.md. Đọc CLAUDE.md, docs/PROGRESS.md và mục "Danh sách issue" của file P05. Làm trên nhánh docs/issues-from-dd, tạo từ origin/main mới nhất.

A. Tạo milestone MVP và các label (mvp, v1, v2, design, bug, windows, wsl) nếu chưa có.

B. Tạo đúng 5 issue theo bảng "Danh sách issue", theo thứ tự trong bảng, gắn milestone MVP và label mvp. Không hỏi lại về cách tách. Body mỗi issue gồm:
- Dòng đầu "Implements DD-xx (slice DD-xxy)" kèm link tới DD.
- Phạm vi: nội dung của slice, chép từ bảng "Implementation slices" trong DD.
- Acceptance criteria: chép đúng từ mục acceptance/test của DD những tiêu chí thuộc slice (theo danh sách TC của slice); không tự thêm bớt. Tiêu chí chung (CI xanh, go test -race) chép vào mọi issue.
- Danh sách package bị ảnh hưởng.
- "Depends on" (số issue) và "Can run in parallel with" (số issue), theo cột Phụ thuộc và Nhóm chạy.

C. Mở một PR docs trên nhánh docs/issues-from-dd:
- Cột Issue trong bảng Traceability (docs/requirements/README.md) và trong docs/basic-design/README.md.
- docs/PROGRESS.md: bước 3 → In progress; thêm dòng Log.
- docs/prompts/README.md: dòng P05 → Done. Không thêm dòng P06…P10; cố vấn thêm khi soạn các file đó.
- go run ./tools/docs-check pass. Mở PR ở dạng draft, chuyển sang ready khi CI xanh. Không merge.

D. Kết thúc:
- Lưu báo cáo thành docs/reports/P05-<YYYY-MM-DD>.md (thêm -2, -3… nếu đã có báo cáo P05 cùng ngày), commit vào PR này và push (P04b-D10).
- In báo cáo theo mẫu CLAUDE.md §13, tiêu đề "Báo cáo — P05 — <YYYY-MM-DD>". Mục 3 liệt kê 5 issue (số, tiêu đề) và chỉ rõ issue nào giao được cho nhiều agent làm song song. Không bắt đầu code.
```

## Sau P05

Mỗi issue được giao bằng một prompt riêng (P04b-D11):

| Prompt | Issue | Chạy khi |
| --- | --- | --- |
| P06 | DD-00a | Ngay sau P05 |
| P07 | DD-00b | Sau khi DD-00a merge; song song với P08, P09 |
| P08 | DD-00c | Sau khi DD-00a merge; song song với P07, P09 |
| P09 | DD-03a | Sau khi DD-00a merge; song song với P07, P08 |
| P10 | DD-03b | Sau khi DD-00b (P07) và DD-03a (P09) merge |

Cố vấn soạn các prompt này dựa trên mẫu C trong `P00-kickoff.md`. Mỗi prompt mở đầu bằng dòng "Implements DD-xx", có mục "Việc Danny làm trước" (bảng đặt file nếu có, và kết quả `git status --short` cần thấy), và kết thúc bằng việc lưu báo cáo vào `docs/reports/` (P04b-D10) rồi in báo cáo theo §13.
