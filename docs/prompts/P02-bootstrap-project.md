---
id: P02
title: Bootstrap the Go project
step: "2a"
status: Ready
branch: chore/bootstrap
pr: ""
depends_on: [P01b]
decisions: [P02-D1, P02-D2, P02-D3, P02-D4]
last_updated: 2026-11-13
---

# P02 — Dựng khung dự án

## Điều kiện chạy

1. PR #3 (P01b) đã merge vào `main`. Sau đó chạy `git checkout main && git pull`.
2. Thay `docs/prompts/P02-bootstrap-project.md` bằng bản file này, không commit; agent sẽ commit trên nhánh mới.
3. Cài golangci-lint: `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`. Không cài cũng được, CI vẫn kiểm tra.
4. Dòng `LICENSE:` trong prompt đang để MIT, là đề xuất của cố vấn. Muốn dùng Apache-2.0 thì sửa dòng đó trước khi dán.

## Quyết định từ review P01b (cố vấn)

- **P02-D1 — Trạng thái prompt.** `Done` nghĩa là phiên đã chạy xong. PR đã merge hay chưa thì xem ở cột Nhánh/PR của mục lục (`(in review)` hoặc `(merged)`), không thêm trạng thái mới. Đây là đề xuất của agent, cố vấn chấp nhận.
- **P02-D2 — Sửa prompt đã chạy.** Chỉ được sửa `status`, `branch` và `pr` trong front matter. Nội dung prompt giữ nguyên.
- **P02-D3 — Giới hạn 200 dòng của docs-check** chỉ áp dụng cho P01, nay bỏ. Khi main.go vượt khoảng 250 dòng, tách thành nhiều file trong cùng package (ví dụ rules.go, frontmatter.go).
- **P02-D4 — Chấp nhận mọi điểm lệch trong báo cáo P01b**, gồm: CL viết bằng tiếng Anh, mẫu §13 đặt trong code block, sửa thêm tham chiếu trong §3 và trong P00. Lần sau, khi tin nhắn của Danny khác với nội dung file prompt, agent làm theo tin nhắn và ghi vào mục 4 của báo cáo, như lần này đã làm.

## Prompt (dán nguyên văn)

```text
Prompt này là docs/prompts/P02-bootstrap-project.md. Đọc CLAUDE.md, docs/PROGRESS.md, docs/prompts/README.md và docs/prompts/P00-kickoff.md. Thực hiện đúng phần B, Bước 0–2 của P00, trên nhánh chore/bootstrap tạo từ origin/main. Mang theo file P02 đã cập nhật mà tôi chép vào (chưa commit).

LICENSE: MIT
(Đây là câu trả lời cho câu hỏi license ở Bước 0 của P00; không cần hỏi lại. Copyright holder: Danny (duong2604), năm 2026.)

Bổ sung so với P00:
- CI: thêm một job chạy `go run ./tools/docs-check`, chỉ cần chạy trên ubuntu.
- Khi tạo go.mod ở gốc repo: gộp tools/docs-check vào module gốc, xóa tools/docs-check/go.mod và go.work (theo P01b-D2), rồi xác nhận `go run ./tools/docs-check` và test của nó vẫn chạy.
- Áp dụng P02-D1…D4 trong file này. Thêm P02-D1 và P02-D2 vào mục "Quy ước" của docs/prompts/README.md. Ghi P02-D3 vào doc comment đầu tools/docs-check/main.go.
- Ghi chú phiên bản: go.mod dùng bản Go đang cài trên máy (1.26.x). CI cài cùng phiên bản đó bằng actions/setup-go với go-version-file: go.mod.

Không tạo issue tính năng, không viết logic tính năng: các phần đó chờ DD được duyệt.

Khi xong:
- Cập nhật docs/PROGRESS.md: bước 1b → Done (PR #3 merged); bước 2a → In review kèm link PR; thêm dòng Log.
- Cập nhật dòng P02 trong docs/prompts/README.md.
- Kết thúc bằng báo cáo theo mẫu CLAUDE.md §13, tiêu đề "Báo cáo — P02 — <YYYY-MM-DD>". Mục 7 ghi rõ kết quả CI trên từng OS (ubuntu, windows, macos).
```
