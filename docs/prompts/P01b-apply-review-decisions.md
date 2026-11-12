---
id: P01b
title: Apply review decisions for P01
step: "1b"
status: Done
branch: docs/apply-p01-review
pr: "#3"
depends_on: [P01]
decisions: [P01b-D1, P01b-D2, P01b-D3, P01b-D4, P01b-D5, P01b-D6, P01b-D7]
clarifications: [CL-01, CL-02, CL-03, CL-04]
last_updated: 2026-11-13
---

# P01b — Áp dụng quyết định review cho P01

Chạy trên nhánh `docs/normalize-structure`, nhánh đã push sau P01.

## Việc Danny làm trước

1. Đăng nhập GitHub CLI: trong Claude Code gõ `! gh auth login` (hoặc chạy `gh auth login` trong terminal).
2. Xem mục B của prompt. Dòng CL nào bạn không đồng ý thì sửa; để nguyên nghĩa là đồng ý.
3. Chép cả thư mục `prompts/` cố vấn gửi vào `docs/prompts/` trên nhánh này, gồm `README.md` và các file `P01…P03`. Không xóa file cũ: agent sẽ dọn ở bước E.

## Prompt (dán nguyên văn)

```text
Tiếp tục trên nhánh docs/normalize-structure. Đọc CLAUDE.md và docs/prompts/README.md. Prompt này là docs/prompts/P01b-apply-review-decisions.md. Áp dụng đúng phạm vi các mục A–E dưới đây; ngoài ra không sửa gì khác. Khi ghi lại quyết định, luôn dùng mã đầy đủ (ví dụ P01b-D5), không dùng mã rút gọn D5.

A. Quyết định kỹ thuật (cố vấn):
P01b-D1 — Một bộ trạng thái cho tài liệu BD, DD, ADR: Draft | In review | Approved | Superseded | Rejected; Rejected chỉ dùng cho ADR. Sửa CLAUDE.md §8: "Status: Proposed" → "status: In review"; "Accepted" → "Approved"; câu "Don't build on a Proposed ADR until it's Accepted" → "Don't build on an ADR until it is Approved". Sửa P00-kickoff.md (sau khi đổi tên ở bước E): "(Accepted)" → "(Approved)". Sửa 3 template và tập trạng thái hợp lệ trong docs-check.
P01b-D2 — Giữ go.work tạm thời. P02 sẽ xóa go.work và tools/docs-check/go.mod khi tạo go.mod ở gốc repo.
P01b-D3 — Chấp nhận hai điểm docs-check làm khác yêu cầu: cho phép .gitkeep, và kiểm tra mã hóa các file *.md ở gốc repo. Ghi hai điểm này vào doc comment đầu main.go và vào CLAUDE.md §3 "Format rules".
P01b-D4 — Chấp nhận ghi đủ 8 phát hiện review.
P01b-D5 — Ngày tháng dùng ISO 8601 YYYY-MM-DD ở mọi nơi: front matter, PROGRESS, ADR, ví dụ trong template (kể cả chỗ "dd-mm-yyyy" trong ADR-template.md). Sửa mọi ngày dạng dd-mm-yyyy hiện có. Thêm quy tắc (e) vào docs-check: last_updated phải khớp ^\d{4}-\d{2}-\d{2}$ và là ngày hợp lệ (time.Parse("2006-01-02")), kèm test bảng.
P01b-D6 — Thêm mục §13 "Reporting & escalation" vào CLAUDE.md, nguyên văn như sau:

## 13. Reporting & escalation

End every session with a report to Danny in Vietnamese, at most ~40 lines, using exactly these headings:

  ## Báo cáo — <prompt ID, e.g. P02> — <YYYY-MM-DD>
  1. Kết quả: Done | Partial | Blocked — one sentence
  2. Nhánh / PR / commit cuối
  3. Đã làm (≤ 5 bullets)
  4. Lệch so với prompt hoặc thiết kế — every deviation, however small: what and why ("Không có" if none)
  5. Câu hỏi kỹ thuật → cố vấn — each with your recommended option
  6. Câu hỏi requirement/spec → Danny — each with the interpretation you used meanwhile
  7. Kiểm chứng — OS, commands run, results; what is NOT verified
  8. Việc Danny cần làm — exact commands

Routing: technical questions (architecture, libraries, git behaviour, file formats, conventions) are decided by the advisor; Danny only reviews. Questions about requirements and product behaviour (what the user sees, defaults, safety policy, naming) are answered by Danny. When unsure which kind, put it under 5. Do not block on a question if a reversible default exists: use it, list it, continue.
Every prompt lives in docs/prompts/ with an ID (see docs/prompts/README.md). Decisions are cited by full ID (e.g. P01b-D5); requirement clarifications by CL-xx. Update docs/PROGRESS.md and docs/prompts/README.md before writing the report.

P01b-D7 — Quy ước đặt tên file prompt, áp dụng ở bước E:
  - Tên file: P<NN>[<chữ cái>]-<kebab-slug>.md. NN là số bước theo thứ tự chạy; chữ cái b, c… dùng cho prompt sửa theo review của chính bước đó (P01 → P01b → P01c). Mỗi prompt mới lấy số kế tiếp, kể cả prompt implement issue (ví dụ P04-implement-dd-00-foundation.md).
  - Mã quyết định: <ID prompt>-D<n> (P01b-D1). Mã làm rõ requirement: CL-xx, đánh số liên tục trên toàn dự án, ghi trong docs/requirements/README.md kèm ID prompt đã chốt nó.
  - Mỗi file prompt có front matter: id, title, step, status (Ready | Running | Done | Superseded), branch, pr, depends_on, decisions, last_updated; thêm clarifications nếu có.
  - Prompt đã chạy thì không sửa nội dung nữa. Muốn thay đổi thì tạo prompt mới với chữ cái kế tiếp, và đổi status prompt cũ thành Superseded nếu nó bị thay thế.
  - docs/prompts/README.md là mục lục: ID | File | Bước | Trạng thái | Nhánh/PR | Quyết định.

B. Câu trả lời của Danny cho câu hỏi requirement. Ghi vào docs/requirements/README.md thành mục mới "## Clarifications", đặt ngay sau "## Original text". KHÔNG sửa phần nguyên văn. Mỗi dòng gồm mã CL-xx, REQ liên quan, câu trả lời, nguồn "P01b", ngày 2026-11-13:
CL-01 (REQ-03) Không lưu "dự án hiện tại". Mỗi lệnh xác định dự án qua project@task, hoặc qua thư mục đang đứng (repo hay worktree của dự án).
CL-02 (REQ-04) Tên nhánh mặc định = <prefix><task>. Prefix cấu hình theo dự án, mặc định "feat/". Có flag --branch <tên đầy đủ> để đặt tên nhánh khác.
CL-03 (REQ-07) Không xóa worktree khi: là worktree chính, có thay đổi chưa commit, có commit chưa push, hoặc đang bị lock. Có --force để bỏ qua; khi chạy trong terminal tương tác thì hỏi xác nhận một lần cho cả danh sách.
CL-04 (4.8) "workspace\" nghĩa là thư mục chứa các dự án (projects root). "nhánh chỉ định chính để tạo workspace" nghĩa là base branch dùng để tạo worktree.
Kèm một dòng ghi chú: các điểm mơ hồ còn lại đã được cố vấn quyết định trong BD-00-overview.md (base lấy từ origin/<base>; đường dẫn worktree có tên dự án; sync = fetch rồi rebase/merge thay cho pull; phần "cần tư vấn" của REQ-06 được trả lời ở BD-00 mục 9).

C. docs-check cho thư mục prompts:
- File trong docs/prompts/, trừ README.md, phải có tên khớp ^P\d{2}[a-z]?-[a-z0-9-]+\.md$.
- Front matter của file prompt phải đủ các trường ở P01b-D7, và status thuộc Ready | Running | Done | Superseded.
- Thêm test bảng cho cả hai quy tắc.

D. Mục lục docs/prompts/README.md (Danny đã chép bản của cố vấn vào): sau khi chạy xong bước E, cập nhật đúng trạng thái và nhánh/PR của từng dòng.

E. Dọn file prompt cũ:
- git mv docs/prompts/KICKOFF_PROMPT.md docs/prompts/P00-kickoff.md và thêm front matter (id: P00, step: "0", status: Ready, decisions: []). Riêng P00 được phép sửa nội dung theo P01b-D1, và đổi câu tham chiếu "NEXT_STEPS_PROMPTS.md, Prompt 3" thành "P03-create-issues-from-dd.md".
- Xóa docs/prompts/NEXT_STEPS_PROMPTS.md (nội dung đã tách sang P01, P01b, P02, P03).
- Sửa mọi tham chiếu tới hai file cũ trong CLAUDE.md, docs/PROGRESS.md, docs/basic-design/README.md và các template.

Hoàn tất:
- Chạy go test trong tools/docs-check và go run ./tools/docs-check: phải sạch lỗi.
- Commit theo nhóm:
    docs: unify status values and ISO dates (P01b-D1, P01b-D5)
    chore(docs-check): validate dates and prompt files
    docs: add reporting protocol to CLAUDE.md (P01b-D6)
    docs: add requirement clarifications CL-01..04
    docs: split prompts into P00–P03 with an index (P01b-D7)
- Push, rồi mở PR "docs: normalize documentation structure" bằng nội dung đã soạn, bổ sung mục "Review decisions P01b-D1…D7, CL-01…04". Không merge.
- Cập nhật docs/PROGRESS.md: bước 1 → In review kèm link PR; thêm phát hiện 9 (status values inconsistent), 10 (date format deviated from template), 11 (go.work temporary), 12 (report format), 13 (prompt files not identifiable per task), mỗi phát hiện ghi trạng thái và mã quyết định xử lý nó; thêm dòng Log.
- Kết thúc bằng báo cáo theo đúng mẫu §13, tiêu đề "Báo cáo — P01b — <ngày>".
```
