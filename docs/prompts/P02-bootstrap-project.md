---
id: P02
title: Bootstrap the Go project
step: "2a"
status: Ready
branch: chore/bootstrap
pr: ""
depends_on: [P01b]
decisions: []
last_updated: 2026-11-13
---

# P02 — Dựng khung dự án

**Điều kiện chạy:** PR của P01/P01b đã merge vào `main`. Cài golangci-lint trên máy trước khi chạy để agent kiểm tra được ở máy; nếu chưa cài thì CI vẫn kiểm tra.

## Prompt (dán nguyên văn)

```text
Prompt này là docs/prompts/P02-bootstrap-project.md. Đọc CLAUDE.md, docs/PROGRESS.md và docs/prompts/P00-kickoff.md. Thực hiện đúng phần B, Bước 0–2 của P00, trên nhánh chore/bootstrap tạo từ origin/main.
Bổ sung so với P00:
- Trong CI, thêm một job chạy `go run ./tools/docs-check` (chỉ cần chạy trên ubuntu).
- Khi tạo go.mod ở gốc repo, gộp tools/docs-check vào module gốc, xóa tools/docs-check/go.mod và go.work (theo P01b-D2), rồi xác nhận `go run ./tools/docs-check` vẫn chạy.
Không tạo issue tính năng và không viết logic tính năng: các phần đó chờ DD được duyệt.
Khi xong: cập nhật docs/PROGRESS.md (bước 2a → In review, thêm link PR, thêm dòng Log) và dòng P02 trong docs/prompts/README.md, rồi kết thúc bằng báo cáo theo mẫu CLAUDE.md §13, tiêu đề "Báo cáo — P02 — <ngày>".
```
