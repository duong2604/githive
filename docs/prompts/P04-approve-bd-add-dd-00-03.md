---
id: P04
title: Approve BD-00 and BD-03, add DD-00, DD-03 and ADR-0004
step: "2d"
status: Done
branch: docs/bd-00-03
pr: "#6"
depends_on: [P03]
decisions: [P04-D1, P04-D2, P04-D3, P04-D4, P04-D5, P04-D6, P04-D7, P04-D8]
clarifications: [CL-05, CL-06]
last_updated: 2026-11-13
---

# P04 — Duyệt BD-00, BD-03; thêm DD-00, DD-03 và ADR-0004

Làm tiếp trên nhánh `docs/bd-00-03`, tức PR #6. Khi xong, PR #6 chứa đủ BD (Approved), DD (Approved) và ADR-0004. Danny merge PR này là bước chốt cuối cùng.

## Việc Danny làm trước

1. `git fetch origin && git switch docs/bd-00-03 && git pull --ff-only`.
2. Chép 6 file cố vấn gửi, không commit:
   - `BD-00-foundation.md`, `BD-03-target.md` → `docs/basic-design/` (ghi đè; bản Approved)
   - `DD-00-foundation.md`, `DD-03-target.md` → `docs/detail-design/`
   - `0004-foundation-contract-additions.md` → `docs/adr/`
   - `P04-approve-bd-add-dd-00-03.md` → `docs/prompts/`
3. Mở hai file BD trên GitHub để xem sơ đồ Mermaid có hiển thị được không, rồi báo cố vấn.

## Câu trả lời của Danny (2026-11-13)

- **CL-05** (NFR-05 / BD-00 BR-13): output, help và thông báo lỗi của CLI chỉ dùng tiếng Anh.
- **CL-06** (REQ-04 / BD-03 BR-04): tên task chỉ dùng chữ thường. Muốn tên nhánh khác, kể cả có chữ hoa, thì dùng `--branch`.

## Quyết định (cố vấn), trả lời mục 5 của báo cáo P03

- **P04-D1 — ADR-0004** ghi lại: exit code 130, JSON envelope và mã lỗi, `GITHIVE_LOCK_TIMEOUT`, việc khóa quanh `git config --local`, file `githive.lock.info`, package `internal/apperr`, dependency `golang.org/x/term`, các flag dùng chung. Commit nguyên văn.
- **P04-D2 — Thông tin người giữ khóa** nằm trong file riêng `githive.lock.info`, theo đề xuất của agent. Đã ghi trong BD-00 BR-10, DD-00 §5.3 và ADR-0004.
- **P04-D3 — Phần code còn lệch BD-00** (envelope, mã lỗi, các flag) được làm theo DD-00. Không sửa code trong P04.
- **P04-D4 — Ví dụ trong BD-03 §3.4** sửa thành `task_18`. Đã có trong bản Approved.
- **P04-D5 — Mục lục BD:** cột "Depends on" của BD-03 ghi `BD-00`. Cột trạng thái BD-00 và BD-03 → Approved; cột DD-00 và DD-03 → Approved, kèm link.
- **P04-D6 — CLAUDE.md §4:**
  - `sync` nhận nhiều target: `githive sync <targets...>`.
  - Danh sách exit code thêm `130 interrupted`.
  - Dòng "Global flags" bỏ chữ "(In review)" và thêm câu: "With --json, stdout holds one envelope {ok, command, data, error, warnings} (BD-00 §3.3, ADR-0004)".
- **P04-D7 — CLAUDE.md §5:**
  - Thêm `internal/apperr/` vào cây thư mục (dòng "errors: Error type, code catalogue, exit codes").
  - Danh sách dependency thêm `golang.org/x/term`.
  - Gạch đầu dòng "Locking" thêm `githive.lock.info` và `GITHIVE_LOCK_TIMEOUT`.
- **P04-D8 — Prompt P05 (chưa chạy, được sửa):**
  - Đổi mọi chữ "P03" trong nội dung thành "P05": tiêu đề, đường dẫn trong prompt, tiêu đề báo cáo, mục "Sau P03".
  - Mục "Sau P05" ghi danh sách issue theo các phần trong DD-00 §1: **DD-00a** trước; sau khi DD-00a merge thì **DD-00b**, **DD-00c**, **DD-03** chạy song song. Các prompt implement sẽ là P06 (DD-00a), P07 (DD-00b), P08 (DD-00c), P09 (DD-03).
  - Bỏ câu "đề xuất cách tách thành nhiều issue và hỏi tôi" cho DD-00, vì DD-00 §1 đã tách sẵn.

## Prompt (dán nguyên văn)

```text
Prompt này là docs/prompts/P04-approve-bd-add-dd-00-03.md. Đọc CLAUDE.md và các phần "Câu trả lời của Danny" và "Quyết định" của file P04. Làm tiếp trên nhánh docs/bd-00-03 (PR #6); mang theo 6 file tôi đã chép vào (chưa commit).

A. Commit nguyên văn 5 file của cố vấn: BD-00-foundation.md, BD-03-target.md (bản Approved), DD-00-foundation.md, DD-03-target.md, docs/adr/0004-foundation-contract-additions.md. Không sửa nội dung; thấy lỗi thì ghi vào mục 5 của báo cáo.

B. Ghi CL-05 và CL-06 vào bảng Clarifications trong docs/requirements/README.md (Source: P04, Date: 2026-11-13), viết bằng tiếng Anh như các dòng CL hiện có.

C. Áp dụng P04-D5…D8 đúng phạm vi mô tả trong file P04. Không sửa code Go.

D. Cập nhật:
- Bảng Traceability: các dòng có BD-00 hoặc BD-03 → "Designed" nếu REQ-format.md có trạng thái đó; nếu không có, giữ "In design" và báo lại trong báo cáo.
- docs/PROGRESS.md: bước 2b, 2c → Done; 2d → In review (PR #6); thêm dòng Log; "Current phase" ghi bước tiếp theo là P05 (tạo issue).
- docs/prompts/README.md: thêm dòng P04; dòng P03 → Done.

E. Rà nhất quán, chỉ báo cáo:
- DD-00, DD-03 so với BD-00, BD-03, ADR-0002, ADR-0003, ADR-0004 và CLAUDE.md §4–§5 sau khi đã sửa.
- Mỗi BR trong BD-00 và BD-03 phải có ít nhất một TC hoặc acceptance criterion trong DD tương ứng. Liệt kê BR nào chưa có.

F. Hoàn tất:
- go run ./tools/docs-check và go test ./... đều pass.
- Commit theo nhóm:
    docs: approve BD-00 and BD-03 (CL-05, CL-06)
    docs: add DD-00 foundation and DD-03 target grammar
    docs: add ADR-0004 foundation contract additions (P04-D1)
    docs: update CLAUDE.md and BD index for approved designs (P04-D5..D7)
    docs: renumber and slice issue prompt P05 (P04-D8)
    docs: record P04 in progress tracker and prompt index
- Push. Cập nhật nội dung PR #6: đổi tiêu đề thành "docs: BD-00, BD-03, DD-00, DD-03 and ADR-0004 (P03, P04)"; thêm mục "Review decisions P04-D1…D8, CL-05, CL-06". Không merge.
- Kết thúc bằng báo cáo theo mẫu CLAUDE.md §13, tiêu đề "Báo cáo — P04 — <YYYY-MM-DD>". Mục 7 ghi kết quả CI.
```
