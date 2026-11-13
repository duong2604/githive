---
id: PROMPTS-INDEX
title: Prompt index
status: Approved
req: []
depends_on: []
owner: Danny (duong2604)
last_updated: 2026-11-14
---

# Prompts

Mỗi phiên Claude Code chạy đúng một prompt. Đến P05, prompt do cố vấn soạn. Từ PO-01 (CLAUDE.md §14): PO gửi `PO-xx` cho TL, TL soạn `Pxx` cho Dev; Danny chạy. Agent cập nhật cột Trạng thái và Nhánh/PR ở cuối mỗi phiên.

## Index

| ID | File | From → To | Bước | Trạng thái | Nhánh / PR | Quyết định |
| --- | --- | --- | --- | --- | --- | --- |
| P00 | [P00-kickoff.md](P00-kickoff.md) | Advisor → Agent | 0 | Ready | — | — |
| P01 | [P01-normalize-docs.md](P01-normalize-docs.md) | Advisor → Agent | 1 | Done | `docs/normalize-structure` / [#1](https://github.com/duong2604/githive/pull/1) (merged) | — |
| P01b | [P01b-apply-review-decisions.md](P01b-apply-review-decisions.md) | Advisor → Agent | 1b | Done | `docs/apply-p01-review` / [#3](https://github.com/duong2604/githive/pull/3) (merged) | P01b-D1…D7, CL-01…04 |
| P01c | [P01c-refresh-bd-00-overview.md](P01c-refresh-bd-00-overview.md) | Advisor → Agent | 1c | Done | `docs/refresh-bd-00` / [#5](https://github.com/duong2604/githive/pull/5) (merged) | P01c-D1…D4 |
| P01d | [P01d-apply-p01c-review.md](P01d-apply-p01c-review.md) | Advisor → Agent | 1d | Done | `docs/refresh-bd-00` / [#5](https://github.com/duong2604/githive/pull/5) (merged) | P01d-D1…D7 |
| P02 | [P02-bootstrap-project.md](P02-bootstrap-project.md) | Advisor → Agent | 2a | Done | `chore/bootstrap` / [#4](https://github.com/duong2604/githive/pull/4) (merged) | P02-D1…D4 |
| P03 | [P03-add-bd-00-and-bd-03.md](P03-add-bd-00-and-bd-03.md) | Advisor → Agent | 2b | Done | `docs/bd-00-03` / [#6](https://github.com/duong2604/githive/pull/6) (merged) | P03-D1…D4 |
| P04 | [P04-approve-bd-add-dd-00-03.md](P04-approve-bd-add-dd-00-03.md) | Advisor → Agent | 2d | Done | `docs/bd-00-03` / [#6](https://github.com/duong2604/githive/pull/6) (merged) | P04-D1…D8, CL-05, CL-06 |
| P04b | [P04b-apply-p04-review.md](P04b-apply-p04-review.md) | Advisor → Agent | 2d | Done | `docs/bd-00-03` / [#6](https://github.com/duong2604/githive/pull/6) (merged) | P04b-D1…D11 |
| P05 | [P05-create-issues-from-dd.md](P05-create-issues-from-dd.md) | Advisor → Agent | 3 | Ready | `docs/issues-from-dd` | — |
| PO-01 | [PO-01-operating-model.md](PO-01-operating-model.md) | PO → TL | 3 | Done | `chore/po-01-operating-model` / [#7](https://github.com/duong2604/githive/pull/7) (merged) | PO-01-D1…D6 |
| PO-02 | [PO-02-ui-spec.md](PO-02-ui-spec.md) | PO → TL | 3 | Done | `docs/po-02-ui-spec` / [#8](https://github.com/duong2604/githive/pull/8) (in review) | PO-02-N1…N4, CL-07…CL-10 |

## Quy ước

- **Tên file:** `P<NN>[<chữ cái>]-<kebab-slug>.md` cho prompt TL → Dev (trước PO-01: cố vấn → agent); `PO-<NN>[<chữ cái>]-<kebab-slug>.md` cho prompt PO → TL (PO-01).
  - `NN` là thứ tự chạy, mỗi prompt mới lấy số kế tiếp; dãy `P` và dãy `PO` đánh số riêng.
  - Prompt từ PO-01 trở đi có thêm `from` và `to` trong front matter (R3, R4).
  - Chữ cái `b`, `c`… đánh dấu prompt sửa theo review của chính bước đó, ví dụ P01 → P01b → P01c.
- **Mã quyết định:** `<ID prompt>-D<n>`, ví dụ `P01b-D5`. Mọi chỗ trích dẫn (commit, PR, PROGRESS, ADR) đều ghi mã đầy đủ.
- **Mã làm rõ requirement:** `CL-xx`, đánh số liên tục trên toàn dự án. Nội dung ghi trong `docs/requirements/README.md`, kèm ID prompt đã chốt nó.
- **Trạng thái prompt:** `Ready` → `Running` → `Done`; `Superseded` khi bị prompt khác thay thế.
- **Prompt đã chạy thì không sửa nội dung.** Muốn thay đổi thì tạo prompt mới với chữ cái kế tiếp.
- **`Done` nghĩa là phiên đã chạy xong (P02-D1).** PR đã merge hay chưa thì xem cột Nhánh / PR (`(in review)` hoặc `(merged)`); không thêm trạng thái mới.
- **Sửa prompt đã chạy (P02-D2):** chỉ được sửa `status`, `branch` và `pr` trong front matter; nội dung giữ nguyên.
- **Báo cáo (P04b-D10):** mỗi phiên lưu báo cáo §13 vào `docs/reports/<ID prompt>-<YYYY-MM-DD>.md` (thêm `-2`, `-3`… nếu cùng prompt chạy nhiều phiên trong một ngày), commit cùng PR trước khi in ra chat; docs-check quy tắc (g) kiểm tên file.
- **Vai trò (PO-01):** quyết định kỹ thuật do TL đưa ra; câu hỏi sản phẩm/UX gửi PO; câu hỏi business/requirement gốc gửi Danny (xem `CLAUDE.md` §13–§15). Quyết định cố vấn trước PO-01 vẫn giữ nguyên giá trị.

## Vòng lặp

1. PO gửi `PO-xx`; Danny dán cho TL.
2. TL thiết kế (DD/ADR) và soạn `Pxx` trong repo (R4); Danny dán cho Dev.
3. Dev làm việc, mở PR draft, giữ file bàn giao (R7) và kết thúc bằng báo cáo theo §13.
4. TL review PR của Dev, rồi báo cáo kèm mục 9 "PO summary"; Danny chuyển cho PO.
5. Danny trả lời câu hỏi business/requirement, duyệt và merge.
