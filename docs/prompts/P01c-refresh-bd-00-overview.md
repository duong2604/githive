---
id: P01c
title: Refresh BD-00 overview to current decisions
step: "1c"
status: Ready
branch: docs/refresh-bd-00
pr: ""
depends_on: [P01b]
decisions: [P01c-D1, P01c-D2, P01c-D3, P01c-D4]
last_updated: 2026-11-13
---

# P01c — Cập nhật BD-00-overview theo các quyết định hiện hành

**Vì sao cần bước này:** BD-00-overview.md được export từ bản thiết kế ngày 2026-11-12, trước khi đổi tên thành GitHive và trước các quyết định trong CLAUDE.md, P01b và CL-01…04. Agent có thể làm theo bản cũ và dựng sai: tên `wtx`, lệnh `new/ls/rm`, khóa cấu hình `wtx.*`, thiếu `project@task`. Cố vấn đã viết lại BD-00 thành bản 2. Agent chỉ thay file và rà lại cho các tài liệu khớp nhau, không tự quyết định thiết kế.

**Thứ tự:** chạy P01c **trước** P02, vì P02 dựng khung code dựa trên tài liệu này.

## Việc Danny làm trước

1. Merge PR #3 nếu chưa merge. Sau đó chạy `git checkout main && git pull`.
2. Chép hai file cố vấn gửi, không commit:
   - `BD-00-overview.md` → `docs/basic-design/BD-00-overview.md` (ghi đè file cũ)
   - `P01c-refresh-bd-00-overview.md` → `docs/prompts/`

## Quyết định (cố vấn)

- **P01c-D1 — BD-00-overview bản 2 là bản chính thức.** Agent commit đúng nguyên văn file Danny đã chép vào. Không sửa nội dung, kể cả chính tả. Thấy lỗi thì ghi vào mục 5 của báo cáo.
- **P01c-D2 — Bỏ bảng ánh xạ tên.** Mục "Naming" trong `docs/basic-design/README.md` (bảng `wtx` → `githive`) không còn cần thiết. Thay bằng một câu: BD-00-overview.md đã dùng tên chính thức từ bản 2 (P01c).
- **P01c-D3 — Thêm `extra_projects` vào cấu hình toàn cục.** Khóa này lưu các dự án đăng ký bằng `githive project add`. Ghi vào CLAUDE.md §5, ở gạch đầu dòng "Config, two layers".
- **P01c-D4 — Phạm vi MVP/v1/v2** theo mục 10 của BD-00 bản 2:
  - Chuyển `--json` và exit code, `copy_files`, `postCreate`, `inspect` lên MVP.
  - Để `shell-init` và phần phát hành (GoReleaser, Scoop/Homebrew, `git-hive`) ở v1.
  - Để `agent-rules` ở v2.

## Prompt (dán nguyên văn)

```text
Prompt này là docs/prompts/P01c-refresh-bd-00-overview.md. Đọc CLAUDE.md, docs/prompts/README.md và file P01c (phần "Quyết định"). Tạo nhánh docs/refresh-bd-00 từ origin/main; mang theo hai file tôi đã chép vào (chưa commit): docs/basic-design/BD-00-overview.md (bản 2) và docs/prompts/P01c-refresh-bd-00-overview.md.

A. Commit BD-00-overview.md bản 2 nguyên văn (P01c-D1).

B. Áp dụng các quyết định:
- P01c-D2: sửa mục "Naming" trong docs/basic-design/README.md.
- P01c-D3: thêm extra_projects vào CLAUDE.md §5.
- P01c-D4: chỉ ghi phạm vi MVP/v1/v2 vào docs/PROGRESS.md (mục "Current phase" hoặc mục mới "Scope"). Không sửa các prompt đã chạy.

C. Rà độ nhất quán trên toàn repo. Đây là bước chỉ kiểm tra, ngoài các chỗ liệt kê ở đây thì không tự sửa thiết kế.
- Tìm mọi chỗ còn dùng tên cũ: grep -rn -i -E "wtx|WTX_|\bnew\b|\bls\b|\brm\b" (lọc kết quả bằng mắt, vì "new", "ls", "rm" hay xuất hiện trong văn bản thường). Bỏ qua phần nguyên văn của prompt đã chạy (P01, P01b), phần lịch sử chỉnh sửa trong BD-00 và văn bản requirement gốc.
- Được tự sửa: tên gọi trong tài liệu hiện hành (CLAUDE.md, docs/basic-design/README.md, docs/PROGRESS.md, docs/templates/, P00-kickoff.md, P02, P03) để khớp với BD-00 bản 2 và CLAUDE.md §4.
- So sánh BD-00 bản 2 với CLAUDE.md §4–§5 theo từng điểm: tên lệnh, alias, cú pháp target, khóa cấu hình, đường dẫn config, file khóa, exit code, cấu trúc thư mục. Điểm nào lệch thì KHÔNG sửa, mà liệt kê vào mục 5 của báo cáo kèm đề xuất.
- Kiểm tra bảng "Basic design index" trong docs/basic-design/README.md có khớp với bảng lệnh ở mục 7 của BD-00 bản 2 không. Lệch thì báo cáo, không sửa.

D. Không sửa P02-bootstrap-project.md: cố vấn đã cập nhật file này (depends_on gồm P01c), Danny sẽ chép vào trước khi chạy P02.

E. Hoàn tất:
- Chạy go run ./tools/docs-check: phải sạch lỗi. Mermaid trong BD-00 không cần kiểm tra bằng công cụ, nhưng hãy mở preview trên GitHub sau khi push và ghi vào báo cáo sơ đồ có hiển thị được không.
- Commit theo nhóm:
    docs: refresh BD-00 overview to GitHive naming (P01c-D1)
    docs: drop wtx naming map, add extra_projects and scope (P01c-D2..D4)
    docs: align current docs with BD-00 rev 2
- Push và mở PR "docs: refresh BD-00 overview (P01c)". Không merge.
- Cập nhật docs/PROGRESS.md (thêm bước 1c, dòng Log) và docs/prompts/README.md (thêm dòng P01c).
- Kết thúc bằng báo cáo theo mẫu CLAUDE.md §13, tiêu đề "Báo cáo — P01c — <YYYY-MM-DD>". Mục 4 liệt kê mọi chỗ đã sửa ở bước C. Mục 5 liệt kê mọi điểm lệch tìm được nhưng không sửa.
```
