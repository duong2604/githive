---
id: P01d
title: Rebase P01c on bootstrap and apply its review decisions
step: "1d"
status: Done
branch: docs/refresh-bd-00
pr: "#5"
depends_on: [P01c, P02]
decisions: [P01d-D1, P01d-D2, P01d-D3, P01d-D4, P01d-D5, P01d-D6, P01d-D7]
last_updated: 2026-11-13
---

# P01d — Rebase PR #5 lên main và áp dụng quyết định review P01c

Bước này làm tiếp trên nhánh `docs/refresh-bd-00` (PR #5), sau khi PR #4 (P02) đã merge.

## Việc Danny làm trước

1. Kiểm tra CI của PR #4: chạy `gh pr checks 4`. Cả 5 job phải xanh: test trên ubuntu, windows, macos; docs-check; goreleaser check. Có job đỏ thì dừng lại, gửi log cho cố vấn.
2. CI xanh thì merge PR #4: `gh pr merge 4 --squash --delete-branch`.
3. Bỏ phần thay đổi đang cất trong stash: `git stash drop stash@{0}`. Đó là ảnh chụp màn hình bị chèn vào P01b và P03, cùng thư mục `docs/prompts/image/`. File prompt không được chứa ảnh; ảnh chụp màn hình thì gửi thẳng cho cố vấn trong chat.
4. Chép file này vào `docs/prompts/`, không commit.

## Quyết định (cố vấn), trả lời mục 5 của báo cáo P01c

- **P01d-D1 — Nghĩa của `--all` và `project@*`.**
  - `githive swarm --all` liệt kê mọi dự án đã đăng ký.
  - `'project@*'` (đặt trong dấu nháy) chọn mọi task của một dự án; dùng được cho `swarm`, `sync` và `retire`.
  - Bỏ cách viết `-p project --all`.
  - Đây là thay đổi cú pháp target, nên phải ghi bằng một ADR mới: `docs/adr/0003-all-flag-and-project-wildcard.md`, status Approved, người quyết định là cố vấn, ghi rõ nó làm rõ (không thay thế) ADR-0002. ADR-0002 giữ nguyên.
- **P01d-D2 — Bảng lệnh trong CLAUDE.md §4** bổ sung ba điểm:
  - `githive version`, có `--json`. Lệnh này đã có từ P02.
  - `githive shell-init`, ghi chú "v1".
  - Thêm `--force` vào `retire`, theo CL-03.
- **P01d-D3 — Gõ `githive` không kèm tham số thì in help**, cho đến khi DD-01 thiết kế menu tương tác. Ghi một câu vào CLAUDE.md §4 và giữ nguyên code.
- **P01d-D4 — Các flag `spawn --branch`, `--from-remote`, `--no-hooks` và `sync --autostash`, `--keep-conflicts`** sẽ được chốt ở DD-04 và DD-06. Không sửa gì ở bước này.
- **P01d-D5 — Quy tắc đặt tên khóa cấu hình.** Khóa trong file TOML dùng snake_case (`post_create`, `copy_files`). Khóa trong git config dùng camelCase, nằm dưới `githive.*` (`githive.postCreate`, `githive.copyFiles`). Hai cách viết này đều đúng, không phải lỗi trong BD-00. Ghi quy tắc vào CLAUDE.md §5.
- **P01d-D6 — Phân biệt `BD-00-overview.md` với dòng BD-00 trong bảng chỉ mục.**
  - Dòng BD-00 trong bảng là tài liệu chức năng "Shared foundation", sẽ có tên `BD-00-foundation.md`, và giữ REQ là NFR-01, 02, 05.
  - BD-00-overview là bản kiến trúc tổng thể nên bao trùm mọi REQ.
  - Thêm một câu giải thích ngay trên bảng trong `docs/basic-design/README.md`. Không đổi front matter của BD-00-overview.
- **P01d-D7 — Ghi chú cho bản phát hành v1.** Khi phát hành, archive cho Windows phải ở dạng `.zip` (dùng `format_overrides` trong GoReleaser). Thêm một dòng vào mục Scope/v1 của `docs/PROGRESS.md`. Chưa sửa `.goreleaser.yaml`.

Kết quả review P02 (PR #4): đạt, nếu CI xanh. Code đã dùng tên `githive` và cấu trúc thư mục khớp với BD-00 bản 2, nên việc chạy P02 trước P01c không gây sai lệch.

## Prompt (dán nguyên văn)

```text
Prompt này là docs/prompts/P01d-apply-p01c-review.md. Đọc CLAUDE.md và phần "Quyết định" của file P01d. PR #4 đã merge vào main.

A. Rebase:
- git fetch origin && git switch docs/refresh-bd-00 && git rebase origin/main.
- Khi gặp conflict ở docs/PROGRESS.md và docs/prompts/README.md: giữ cả các dòng của P02 (bước 2a) lẫn của P01c (bước 1c). Mục "Current phase" ghi: Phase 1 xong phần nền tảng tài liệu và khung code; bước tiếp theo là 2b (BD-00-foundation, BD-03 → DD-00, DD-03).
- Commit file P01d (tôi đã chép vào, chưa commit).

B. Áp dụng P01d-D1…D7 đúng phạm vi mô tả trong file. Không sửa BD-00-overview.md (vẫn theo P01c-D1), ADR-0002, hay các prompt đã chạy, ngoài front matter được phép theo P02-D2.

C. Hoàn tất:
- go run ./tools/docs-check và go test ./... đều phải pass.
- Commit theo nhóm:
    docs: add ADR-0003 for --all and project@* (P01d-D1)
    docs: update CLAUDE.md command table and config naming (P01d-D2, D3, D5)
    docs: clarify BD-00 overview vs foundation, note v1 zip archives (P01d-D6, D7)
- git push --force-with-lease (chỉ trên nhánh này). Cập nhật nội dung PR #5: thêm mục "Review decisions P01d-D1…D7". Không merge.
- Cập nhật docs/PROGRESS.md (thêm bước 1d, dòng Log) và docs/prompts/README.md (dòng P01d; dòng P01c, P02 ghi đúng PR và trạng thái merge).
- Kết thúc bằng báo cáo theo mẫu CLAUDE.md §13, tiêu đề "Báo cáo — P01d — <YYYY-MM-DD>". Mục 7 ghi kết quả CI của PR #5.
```
