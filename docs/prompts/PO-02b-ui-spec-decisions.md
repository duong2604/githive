---
id: PO-02b
title: PO decisions on the PO-02 report; UI spec Approved
from: PO
to: TL
step: "3"
status: Running
branch: docs/po-02-ui-spec
pr: "#8"
depends_on: [PO-02]
decisions: [UI-D1, UI-D2, UI-D3, UI-D4, UI-D5, UI-D6, UI-D7, UI-D8, UI-D9, UI-D10]
last_updated: 2026-11-14
---

# PO-02b — Quyết định của PO cho báo cáo PO-02; UI spec chuyển sang Approved

Nguyên văn prompt Danny dán vào phiên (R3). Bản dán đã mất định dạng Markdown (xuống dòng, dấu `` ` `` và `*`; ví dụ `'kirindo@*'` thành `'kirindo@'`); văn bản dưới đây giữ đúng như bản dán. Khi áp dụng, TL đọc theo ý định gốc (ví dụ E7 là `'kirindo@*'`) và ghi lại trong báo cáo `docs/reports/PO-02b-2026-11-14.md`.

## Prompt (nguyên văn)

```text
PO-02b — Quyết định của PO cho báo cáo PO-02; UI spec chuyển sang Approved Người gửi: PO · Người nhận: TL · Làm tiếp trên PR #8 (nhánh docs/po-02-ui-spec), trước khi Danny merge.

================================================================== NHẬN XÉT CỦA PO VỀ PO-02
Bạn khôi phục định dạng Markdown của bản dán bị mất là đúng; câu chữ giữ nguyên. Sửa thêm một chỗ: dòng UI-06 và UI-07 trong bảng §3 bị mất dấu ` quanh tên lệnh (xem E1).
Phần đánh giá kỹ thuật (a)–(d) được chấp nhận. Mục a: 1–9, 14–16, 19 là việc của TL, áp dụng ở DD rev 4 trong phiên P05 (PO-02-N4). Mục a: 10–13, 17, 18, 20 thuộc sản phẩm; PO quyết định ở dưới.
Kế hoạch BD-01 (setup, gate, menu) và DD-01a/b/c: đồng ý. Menu làm cuối cùng trong MVP. BD-01 do PO viết, sau khi DD-00 rev 4 xong.
================================================================== QUYẾT ĐỊNH CỦA PO (mã UI-D1…UI-D10; trả lời Q1–Q8 và mục a)

UI-D1 (Q1, a17) MVP không có --merged/--stale. Ví dụ retire trong spec đổi sang cú pháp MVP; hai filter đó thuộc v1. UI-D2 (Q2) Dùng một chữ NULL ở mọi nơi (picker, project list, JSON), theo REQ-08. UI-D3 (Q3) STATE theo thứ tự ưu tiên: locked > dirty > unpushed > merged > clean. external thay cho mọi giá trị khác. inspect hiện tất cả các trạng thái đang đúng. UI-D4 (Q4) retire thoát 0 khi mọi task được chọn đều đã xóa. Thoát 4 (unsafe) khi có ít nhất một task được chọn bị giữ lại do kiểm tra an toàn. Task bị filter loại ra (v1) không tính. JSON liệt kê removed và kept, kept có kèm lý do. UI-D5 (Q5) Output không phải terminal (pipe, file) dùng chữ ASCII (ok, warning:, error:). UI-D6 (Q6) Bấm Esc hoặc q ở picker khi đang chạy một lệnh (không phải trong menu) nghĩa là hủy: in "cancelled", thoát 130 (interrupted). UI-D7 (Q7) Giữ lệnh completion của cobra từ MVP. Gợi ý tên project và task khi bấm Tab thuộc v1. UI-D8 (Q8, a10) Wizard MVP có 4 bước, không hỏi sync strategy (giữ giá trị mặc định cho đến v1). Sửa CLAUDE.md §4: bỏ "sync strategy" khỏi mô tả init. UI-D9 (a20) Các lệnh chạy được khi chưa setup (gate không chặn): init, doctor, version, help, completion. UI-D10 (c, Git Bash) Nếu cửa sổ Git Bash độc lập (mintty) không được nhận là terminal: GitHive chạy ở chế độ plain và in một dòng gợi ý dùng Windows Terminal, terminal của VS Code, hoặc "winpty githive". Ghi là giới hạn đã biết; xác nhận lại sau khi xong DD-00a. Mục a11, a12: spec là chuẩn. Ví dụ trong overview chỉ để minh họa (E10). Mục a13, a18: sửa overview (E11, E12). Mục a15: inspect lấy ngày tạo từ reflog, hiện "unknown" khi không có; chấp nhận.

================================================================== VIỆC CỦA TL
Làm trên nhánh docs/po-02-ui-spec (PR #8): git switch docs/po-02-ui-spec && git pull. Lưu prompt này thành docs/prompts/PO-02b-ui-spec-decisions.md (R3); nếu bản dán mất định dạng, khôi phục như ở PO-02 và ghi vào báo cáo. Cập nhật docs/handoff/TL.md (R7).
Sửa docs/product/ui-spec.md theo từng mục (không viết lại cả file): E1 §3 bảng: dòng UI-06 → "Command results: spawn, swarm, inspect, retire, project list, path, doctor"; dòng UI-07 → "sync screens". E2 §5 Rules: thêm dòng "- The MVP wizard does not ask for a sync strategy; the default is kept until sync arrives in v1 (UI-D8)." E3 §6: thêm cuối mục "Commands that work before setup, without the gate: init, doctor, version, help, completion (UI-D9)." E4 §7 Rules: thêm dòng "- In a command outside the menu, Esc or q in a picker cancels: GitHive prints cancelled and exits 130 (interrupted) (UI-D6)." E5 §8 project picker: "legacy base: NOT SET" → "legacy base: NULL"; thêm câu dưới khối: "NULL is the only word for a missing base branch, everywhere (UI-D2)." E6 §9 swarm: thêm dòng sau dòng "STATE values": "- When several states apply, the first in this order is shown: locked, dirty, unpushed, merged, clean; external replaces all others. inspect lists every state that applies (UI-D3)." E7 §9 retire:
Lệnh ví dụ "$ githive retire 'kirindo@' --merged" → "$ githive retire 'kirindo@'". Giữ nguyên bảng.
Thêm sau câu "The default answer is No…": "Filters --merged and --stale arrive in v1 (UI-D1). Exit code: 0 when every selected task was removed; 4 (unsafe) when at least one selected task was kept by a safety check; tasks left out by a filter do not count. JSON lists removed and kept, each kept task with its reason (UI-D4)." E8 §12: thêm dòng "- Output that is not a terminal (pipe, file) uses the ASCII words (UI-D5)." Thêm dòng "- Known limitation: a standalone Git Bash window (mintty) may not be detected as a terminal; GitHive then runs in plain mode and prints a one-line tip: use Windows Terminal, the VS Code terminal, or winpty githive (UI-D10)." E9 §15: thêm câu "The completion command itself is available from the MVP (UI-D7)." E10 §2 Design principles: thêm dòng cuối "- P6 — This document is the reference for what users see. Output examples in BD and DD documents are illustrative." E11 Thêm mục mới "## 19. PO decisions (PO-02b)" trước "## 18. Traceability" nếu cần giữ thứ tự số, hoặc ngay sau §17; liệt kê UI-D1…UI-D10 mỗi dòng một câu (chép từ phần QUYẾT ĐỊNH CỦA PO, viết bằng tiếng Anh). E12 Front matter: status: Approved; last_updated: 2026-11-14. Thêm cuối file mục "## Revision history" với hai dòng: "| 2026-11-13 | In review (PO-02) |" và "| 2026-11-14 | Approved (PO-02b): UI-D1…UI-D10 |".
Sửa docs/basic-design/BD-00-overview.md (tài liệu này PO sở hữu): O1 §7, dòng githive project list: "Bảng gồm tên dự án, base branch (hoặc NULL), số worktree, remote" → "Bảng gồm tên dự án, base branch (hoặc NULL), số task, đường dẫn (UI spec §9)". O2 §7, ngay trước dòng "Ví dụ output của githive swarm 'kirindo@*':": thêm câu "Màn hình, bảng và thông báo cụ thể theo docs/product/ui-spec.md; các ví dụ output trong tài liệu này chỉ để minh họa." O3 §10, bảng rủi ro, dòng "Tốn dung lượng": thay "swarm --size, retire --stale" bằng "retire hiện dung lượng sẽ giải phóng (CL-08), retire --stale (v1)". O4 §10, roadmap v1: thay "phát hành bằng GoReleaser kèm Scoop/Homebrew và alias git-hive" bằng "phát hành bằng GoReleaser kèm winget, Scoop, Homebrew, script cài đặt và alias git-hive". O5 Lịch sử chỉnh sửa: thêm dòng "| 2026-11-14 | Bản 4 (PO-02b): khớp với UI spec — project list, dung lượng ở retire, kênh phát hành v1 |".
CLAUDE.md §4: bỏ "sync strategy" khỏi mô tả init (UI-D8).
Hoàn tất:
docs-check, go vet, go test đều pass.
Commit theo nhóm: docs: save PO-02b prompt docs: UI spec approved with PO decisions UI-D1..D10 docs: align BD-00 overview and CLAUDE.md with the UI spec docs: record PO-02b
PROGRESS và mục lục prompt: thêm PO-02b.
Chuyển PR #8 sang ready khi CI xanh; không merge.
Lưu báo cáo docs/reports/PO-02b-<ngày>.md. Báo cáo theo §13, có mục 9 "PO summary".
Phiên tiếp theo của bạn: P05 theo PO-02-N4 (DD-00 rev 4, DD-03 rev 4 gồm cả mục a: 1–9, 14–16, 19 và UI-D2…D6, D9, D10 nếu liên quan; rồi 5 issue và P06).
```
