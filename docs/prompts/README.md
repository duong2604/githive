---
id: PROMPTS-INDEX
title: Prompt index
status: Approved
req: []
depends_on: []
owner: Danny (duong2604)
last_updated: 2026-11-13
---

# Prompts

Mỗi phiên Claude Code chạy đúng một prompt. Các prompt dưới đây do cố vấn soạn và Danny chạy. Agent cập nhật cột Trạng thái và Nhánh/PR ở cuối mỗi phiên.

## Index

| ID | File | Bước | Trạng thái | Nhánh / PR | Quyết định |
| --- | --- | --- | --- | --- | --- |
| P00 | [P00-kickoff.md](P00-kickoff.md) | 0 | Ready | — | — |
| P01 | [P01-normalize-docs.md](P01-normalize-docs.md) | 1 | Done | `docs/normalize-structure` / [#1](https://github.com/duong2604/githive/pull/1) (merged) | — |
| P01b | [P01b-apply-review-decisions.md](P01b-apply-review-decisions.md) | 1b | Done | `docs/apply-p01-review` / [#3](https://github.com/duong2604/githive/pull/3) (merged) | P01b-D1…D7, CL-01…04 |
| P01c | [P01c-refresh-bd-00-overview.md](P01c-refresh-bd-00-overview.md) | 1c | Done | `docs/refresh-bd-00` / [#5](https://github.com/duong2604/githive/pull/5) (merged) | P01c-D1…D4 |
| P01d | [P01d-apply-p01c-review.md](P01d-apply-p01c-review.md) | 1d | Done | `docs/refresh-bd-00` / [#5](https://github.com/duong2604/githive/pull/5) (merged) | P01d-D1…D7 |
| P02 | [P02-bootstrap-project.md](P02-bootstrap-project.md) | 2a | Done | `chore/bootstrap` / [#4](https://github.com/duong2604/githive/pull/4) (merged) | P02-D1…D4 |
| P03 | [P03-add-bd-00-and-bd-03.md](P03-add-bd-00-and-bd-03.md) | 2b | Done | `docs/bd-00-03` / [#6](https://github.com/duong2604/githive/pull/6) (in review) | P03-D1…D4 |
| P04 | [P04-approve-bd-add-dd-00-03.md](P04-approve-bd-add-dd-00-03.md) | 2d | Done | `docs/bd-00-03` / [#6](https://github.com/duong2604/githive/pull/6) (in review) | P04-D1…D8, CL-05, CL-06 |
| P05 | [P05-create-issues-from-dd.md](P05-create-issues-from-dd.md) | 3 | Ready | `docs/issues-from-dd` | — |

## Quy ước

- **Tên file:** `P<NN>[<chữ cái>]-<kebab-slug>.md`.
  - `NN` là thứ tự chạy, mỗi prompt mới lấy số kế tiếp.
  - Chữ cái `b`, `c`… đánh dấu prompt sửa theo review của chính bước đó, ví dụ P01 → P01b → P01c.
- **Mã quyết định:** `<ID prompt>-D<n>`, ví dụ `P01b-D5`. Mọi chỗ trích dẫn (commit, PR, PROGRESS, ADR) đều ghi mã đầy đủ.
- **Mã làm rõ requirement:** `CL-xx`, đánh số liên tục trên toàn dự án. Nội dung ghi trong `docs/requirements/README.md`, kèm ID prompt đã chốt nó.
- **Trạng thái prompt:** `Ready` → `Running` → `Done`; `Superseded` khi bị prompt khác thay thế.
- **Prompt đã chạy thì không sửa nội dung.** Muốn thay đổi thì tạo prompt mới với chữ cái kế tiếp.
- **`Done` nghĩa là phiên đã chạy xong (P02-D1).** PR đã merge hay chưa thì xem cột Nhánh / PR (`(in review)` hoặc `(merged)`); không thêm trạng thái mới.
- **Sửa prompt đã chạy (P02-D2):** chỉ được sửa `status`, `branch` và `pr` trong front matter; nội dung giữ nguyên.
- **Vai trò:** quyết định kỹ thuật do cố vấn đưa ra, câu hỏi requirement do Danny trả lời (xem `CLAUDE.md` §13).

## Vòng lặp

1. Danny chạy prompt.
2. Agent làm việc, push nhánh và kết thúc bằng báo cáo theo §13.
3. Danny gửi báo cáo và tên nhánh cho cố vấn.
4. Cố vấn review nhánh, rồi soạn prompt sửa (`Pxxb`) hoặc prompt bước kế tiếp.
5. Danny trả lời câu hỏi requirement, duyệt quyết định và merge.
