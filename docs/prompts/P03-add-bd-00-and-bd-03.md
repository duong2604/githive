---
id: P03
title: Add BD-00-foundation and BD-03 for review; apply P01d follow-ups
step: "2b"
status: Done
branch: docs/bd-00-03
pr: "https://github.com/duong2604/githive/pull/6"
depends_on: [P01d]
decisions: [P03-D1, P03-D2, P03-D3, P03-D4]
last_updated: 2026-11-13
---

# P03 — Đưa BD-00-foundation và BD-03 vào repo để review

Cố vấn đã viết hai tài liệu basic design đầu tiên. Agent commit nguyên văn ở trạng thái **In review**, áp dụng các quyết định còn mở từ báo cáo P01d, và dọn mục lục prompt. Danny review hai tài liệu trên PR này. Khi Danny duyệt, P04 sẽ chuyển chúng sang Approved và thêm DD-00, DD-03.

## Việc Danny làm trước

1. Kiểm tra máy: đang ở `main` mới nhất (có 62dd4b1 của PR #5) và không có thay đổi chưa commit.
2. Chép các file cố vấn gửi vào repo, không commit:
   - `BD-00-foundation.md` → `docs/basic-design/`
   - `BD-03-target.md` → `docs/basic-design/`
   - `BD-00-overview.md` (bản 3) → `docs/basic-design/`, ghi đè bản cũ
   - `P03-add-bd-00-and-bd-03.md` → `docs/prompts/`

## Quyết định (cố vấn)

- **P03-D1 — Bỏ hẳn `-p/--project` ở mọi lệnh.** Dự án được xác định qua target (`project@task`, `'project@*'`) hoặc qua thư mục đang đứng (CL-01, BD-03). Đây là câu trả lời cho câu hỏi 1 trong báo cáo P01d. BD-00-overview bản 3 đã sửa theo quyết định này. ADR-0003 để điểm này cho BD-05 quyết định; nay quyết định sớm hơn ở BD-03. ADR-0003 đã Approved nên không sửa; quyết định được ghi vào CLAUDE.md §4 (xem bước B của prompt).
- **P03-D2 — Đồng bộ bảng lệnh trong `README.md` ở gốc repo** cho giống hệt bảng trong CLAUDE.md §4, gồm cả `shell-init` (v1), `retire --force`, `swarm ['project@*' | --all]`. Câu dưới bảng ghi đúng ngữ pháp target theo BD-03.
- **P03-D3 — Chấp nhận ADR-0002 được đọc kèm ADR-0003.** Không sửa ADR-0002.
- **P03-D4 — Đánh lại số prompt** theo thứ tự chạy thực tế:
  - `P03-create-issues-from-dd.md` (chưa chạy) đổi tên thành `P05-create-issues-from-dd.md`, `id: P05`, `step: "3"`, `depends_on: [P04]`.
  - Trong nội dung file đó, đổi "P04-implement-dd-00-foundation.md và P05-implement-dd-03-target.md" thành "P06-implement-dd-00-foundation.md và P07-implement-dd-03-target.md".
  - P04 (thêm DD-00, DD-03) do cố vấn soạn sau khi Danny duyệt BD.

## Prompt (dán nguyên văn)

```text
Prompt này là docs/prompts/P03-add-bd-00-and-bd-03.md. Đọc CLAUDE.md và phần "Quyết định" của file P03. Tạo nhánh docs/bd-00-03 từ origin/main; mang theo các file tôi đã chép vào (chưa commit): docs/basic-design/BD-00-foundation.md, docs/basic-design/BD-03-target.md, docs/basic-design/BD-00-overview.md (bản 3), docs/prompts/P03-add-bd-00-and-bd-03.md.

A. Commit nguyên văn ba file BD do cố vấn viết. Không sửa nội dung, kể cả chính tả. Thấy lỗi hoặc mâu thuẫn thì ghi vào mục 5 của báo cáo.

B. Áp dụng P03-D1…D4 đúng phạm vi mô tả trong file P03. Riêng CLAUDE.md §4:
- thêm câu "`-p/--project` does not exist; the project comes from the target or the current directory (BD-03, P03-D1)."
- thêm dòng "Global flags: --json, --yes/-y, --verbose/-v, --no-color, --help (BD-00 §3.1, In review)."
Không đổi exit code, error code hay JSON envelope trong code: các phần đó chờ BD-00 được duyệt và DD-00.

C. Cập nhật:
- docs/basic-design/README.md: cột "BD status" của BD-00 và BD-03 → In review; thêm link tới hai file.
- docs/requirements/README.md, bảng Traceability: REQ-03, REQ-04, REQ-05, REQ-07, NFR-01, NFR-02, NFR-05 → "In design" ở những dòng có BD-00 hoặc BD-03 (theo quy tắc trong REQ-format.md).
- docs/PROGRESS.md: bước 1c và 1d → Done (PR #5 merged); bước 2b → In review kèm link PR; thêm dòng Log.
- docs/prompts/README.md: P01c, P01d → Done (merged); thêm P03 (dòng mới); đổi dòng P03 cũ thành P05.

D. Rà nhất quán, chỉ báo cáo, không tự sửa:
- So BD-00-foundation và BD-03 với CLAUDE.md §4–§5, ADR-0002, ADR-0003 và code hiện có (internal/cli). Liệt kê mọi điểm lệch, ví dụ: exit code 130, error code mới, JSON envelope của version, flag --verbose/--no-color chưa có trong code.
- Kiểm tra Mermaid trong hai BD mới hiển thị được trên GitHub sau khi push.

E. Hoàn tất:
- go run ./tools/docs-check và go test ./... đều pass.
- Commit theo nhóm:
    docs: add BD-00 foundation and BD-03 target grammar (In review)
    docs: BD-00 overview rev 3, drop -p/--project (P03-D1)
    docs: sync README command table with CLAUDE.md (P03-D2)
    docs: renumber issue prompt to P05 (P03-D4)
    docs: record P03 in progress tracker and prompt index
- Push, mở PR "docs: add BD-00 foundation and BD-03 for review (P03)". Body ghi rõ: "Danny reviews BD-00-foundation and BD-03 in this PR; open questions are in §9 of each file." Không merge.
- Kết thúc bằng báo cáo theo mẫu CLAUDE.md §13, tiêu đề "Báo cáo — P03 — <YYYY-MM-DD>".
```
