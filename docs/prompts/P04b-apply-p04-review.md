---
id: P04b
title: Apply P04 review decisions before merging PR #6
step: "2d"
status: Done
branch: docs/bd-00-03
pr: "#6"
depends_on: [P04]
decisions: [P04b-D1, P04b-D2, P04b-D3, P04b-D4, P04b-D5, P04b-D6, P04b-D7, P04b-D8, P04b-D9, P04b-D10, P04b-D11]
last_updated: 2026-11-13
---

# P04b — Áp dụng quyết định review P04 trước khi merge PR #6

Báo cáo P04 tái lập (laptop nhà, 2026-11-13) tìm ra 7 điểm trong mục 5 và 3 điểm lệch trong mục 4. Cố vấn đã sửa 4 tài liệu thiết kế (bản rev 2, vẫn Approved). Agent commit các tài liệu đó, sửa vài chỗ nhỏ, và bắt đầu lưu báo cáo vào repo. Sau P04b, cố vấn review lần cuối rồi Danny merge PR #6.

## Việc Danny làm trước

**1. Chuẩn bị máy** (laptop Windows ở nhà, gốc repo `D:\workspace\githive`):

```bash
git switch docs/bd-00-03 && git pull --ff-only
git stash drop stash@{0}      # stash cũ chỉ còn ảnh rác, không cần giữ
```

**2. Đặt file:**

| # | File bạn tải về | Đặt vào (tính từ gốc repo) | Ghi chú |
| --- | --- | --- | --- |
| 1 | `BD-00-foundation.md` | `docs/basic-design/BD-00-foundation.md` | Ghi đè (rev 2) |
| 2 | `BD-03-target.md` | `docs/basic-design/BD-03-target.md` | Ghi đè (rev 2) |
| 3 | `DD-00-foundation.md` | `docs/detail-design/DD-00-foundation.md` | Ghi đè (rev 2) |
| 4 | `DD-03-target.md` | `docs/detail-design/DD-03-target.md` | Ghi đè (rev 2) |
| 5 | `P04-2026-11-13.md` | `docs/reports/P04-2026-11-13.md` | Thư mục mới `docs/reports/` |
| 6 | `P04b-apply-p04-review.md` | `docs/prompts/P04b-apply-p04-review.md` | File mới |

PowerShell (chạy ở gốc repo; sửa `$DL` nếu bạn tải file về chỗ khác):

```powershell
$DL = "$env:USERPROFILE\Downloads"
New-Item -ItemType Directory -Force docs\reports | Out-Null
Copy-Item "$DL\BD-00-foundation.md" docs\basic-design\ -Force
Copy-Item "$DL\BD-03-target.md"     docs\basic-design\ -Force
Copy-Item "$DL\DD-00-foundation.md" docs\detail-design\ -Force
Copy-Item "$DL\DD-03-target.md"     docs\detail-design\ -Force
Copy-Item "$DL\P04-2026-11-13.md"   docs\reports\
Copy-Item "$DL\P04b-apply-p04-review.md" docs\prompts\
git status --short
```

Git Bash (chạy ở gốc repo):

```bash
DL=/c/Users/duong/Downloads
mkdir -p docs/reports
cp "$DL/BD-00-foundation.md" "$DL/BD-03-target.md" docs/basic-design/
cp "$DL/DD-00-foundation.md" "$DL/DD-03-target.md" docs/detail-design/
cp "$DL/P04-2026-11-13.md" docs/reports/
cp "$DL/P04b-apply-p04-review.md" docs/prompts/
git status --short
```

**3. Kết quả `git status --short` phải đúng như sau** (6 dòng, không thêm dòng nào):

```
 M docs/basic-design/BD-00-foundation.md
 M docs/basic-design/BD-03-target.md
 M docs/detail-design/DD-00-foundation.md
 M docs/detail-design/DD-03-target.md
?? docs/prompts/P04b-apply-p04-review.md
?? docs/reports/
```

**4.** Mở Claude Code ở gốc repo, dán phần "Prompt (dán nguyên văn)" bên dưới.

## Quyết định (cố vấn), trả lời báo cáo P04 tái lập

| Mã | Điểm trong báo cáo | Quyết định | Đã sửa ở |
| --- | --- | --- | --- |
| P04b-D1 | 5a: DD-03 cần runner của DD-00b | Đồng ý với agent. Tách DD-03 thành **DD-03a** (parser + resolver, test dùng dữ liệu giả, chạy song song với DD-00b/00c) và **DD-03b** (`gitLocator` + TC-14…18, chạy sau DD-00b). | DD-03 rev 2 §1, DD-00 rev 2 §1 |
| P04b-D2 | 5b, 5c: `platform/wsl.go` chưa thuộc phần nào; `platform` import `ui` | Đồng ý. `internal/platform` (term, wsl, paths) chuyển vào **DD-00a**. Bỏ `RepoWarnings`; `platform` chỉ trả `UnderMnt(path) bool`; cảnh báo `repo_on_mnt` do `target` tạo. `platform` không import `ui` hay `apperr`. | DD-00 rev 2 §1, §2, §3.6; DD-03 rev 2 §5.4 |
| P04b-D3 | 5d: BD-03 và DD-03 lệch về chế độ tương tác | **Khác đề xuất của agent:** sửa theo hướng DD-03 làm đủ những gì BD-03 §5 yêu cầu, tức hỏi chọn cả dự án lẫn task khi thiếu. Lý do: NFR-05 (chọn thay vì gõ lệnh). Bỏ câu mâu thuẫn ở BD-03 §8. Thêm TC-21. | BD-03 rev 2; DD-03 rev 2 §5.4, TC-21 |
| P04b-D4 | 5e: `SpecSwarm` chưa chặn `project@task` | Thêm `Spec.AllowTask`; `swarm` chỉ nhận `'project@*'` hoặc `--all`. Thêm TC-22. | DD-03 rev 2 §3.4, §5.4, TC-22 |
| P04b-D5 | 5f: bảng mã lỗi BD-00 §3.4 thiếu `interrupted` | Thêm dòng 130 `interrupted`. | BD-00 rev 2 |
| P04b-D6 | 5g: BR chưa có TC | DD-00 thêm TC-27 (`NativeAbs`, BR-11) và TC-28 (`lock.With` nhả khóa khi bị hủy, BR-12). BR-13 là quy ước review (CL-05), không cần TC. BD-03 BR-05, BR-09, BR-11 ghi rõ để cho DD-02, DD-05 và tài liệu. | DD-00 rev 2 §9–§11; DD-03 rev 2 §11 |
| P04b-D7 | 4: `README.md` còn `sync <target>` | Sửa thành `sync <targets...>` cho giống CLAUDE.md §4. | Agent sửa |
| P04b-D8 | 4: CLAUDE.md §5 "Locking" thiếu `git config --local` | Thêm vào danh sách thao tác được khóa (ADR-0004 #4). Lỗi của cố vấn ở P04-D7. | Agent sửa |
| P04b-D9 | 4: mục "Not changed on purpose" trong PR #6 đã cũ | Cập nhật nội dung PR #6 cho đúng hiện trạng. | Agent sửa |
| P04b-D10 | Báo cáo bị mất khi đổi máy | Từ nay mỗi phiên lưu báo cáo vào `docs/reports/<ID prompt>-<YYYY-MM-DD>.md` (lần thứ hai trong ngày thêm `-2`), commit cùng PR, trước khi in ra chat. Ghi vào CLAUDE.md §13; docs-check thêm quy tắc tên file. | Agent sửa |
| P04b-D11 | Thứ tự prompt implement | P05 tạo 5 issue: **DD-00a**; sau đó **DD-00b, DD-00c, DD-03a** chạy song song; **DD-03b** sau DD-00b. Prompt implement: P06 (00a), P07 (00b), P08 (00c), P09 (03a), P10 (03b). | Agent sửa P05 |

Ngoài ra: Traceability giữ "In design" vì REQ-format chưa có trạng thái "Designed" — chấp nhận, sẽ thêm trạng thái này khi có DD đầu tiên được code xong.

## Prompt (dán nguyên văn)

```text
Prompt này là docs/prompts/P04b-apply-p04-review.md. Đọc CLAUDE.md và bảng "Quyết định" của file P04b. Làm tiếp trên nhánh docs/bd-00-03 (PR #6); mang theo các file tôi đã chép vào (chưa commit): 4 tài liệu BD/DD rev 2, docs/reports/P04-2026-11-13.md, docs/prompts/P04b-apply-p04-review.md.

A. Commit nguyên văn 4 tài liệu rev 2 (BD-00, BD-03, DD-00, DD-03) và báo cáo P04 tái lập. Không sửa nội dung; thấy lỗi thì ghi vào mục 5 của báo cáo.

B. Áp dụng P04b-D7…D11 đúng phạm vi:
- D7: README.md, dòng sync → `githive sync <targets...>`, giống CLAUDE.md §4.
- D8: CLAUDE.md §5 "Locking": "held only around fetch / worktree add / remove / branch delete / `git config --local` writes".
- D9: sửa nội dung PR #6 (gh pr edit 6 --body-file …): bỏ hoặc cập nhật đoạn "Not changed on purpose"; thêm mục "Review decisions P04b-D1…D11".
- D10:
  - CLAUDE.md §13 thêm đoạn: "Before printing the report, save it as docs/reports/<prompt ID>-<YYYY-MM-DD>.md (add -2, -3… for more sessions of the same prompt on one day) and commit it with the session's other changes, so any machine can read it after git pull."
  - tools/docs-check thêm quy tắc (g): file trong docs/reports/ phải có tên khớp ^P\d{2}[a-z]?-\d{4}-\d{2}-\d{2}(-\d+)?\.md$ và ngày hợp lệ; không yêu cầu front matter. Thêm test bảng.
- D11: sửa docs/prompts/P05-create-issues-from-dd.md (prompt chưa chạy, được sửa):
  - danh sách issue theo đúng P04b-D11 (5 issue, ghi rõ phụ thuộc và nhóm chạy song song);
  - tên prompt implement P06…P10;
  - prompt P05 kết thúc bằng việc lưu báo cáo theo D10;
  - mục "Việc Danny làm trước" của P05 phải có bảng đặt file (nếu có file cần chép) và kết quả git status cần thấy.

C. Cập nhật:
- docs/PROGRESS.md: thêm bước 2d-b (P04b) → In review (PR #6); dòng Log ghi rõ báo cáo P04 gốc bị mất và đã tái lập ở docs/reports/P04-2026-11-13.md.
- docs/prompts/README.md: thêm dòng P04b; mục "Quy ước" thêm một dòng về docs/reports (P04b-D10).

D. Rà nhất quán, chỉ báo cáo:
- DD-00 và DD-03 rev 2 với BD-00, BD-03 rev 2, ADR-0002/0003/0004 và CLAUDE.md (sau B).
- Mọi BR của BD-00, BD-03 đều có TC, tiêu chí nghiệm thu, hoặc ghi chú "để cho DD sau".

E. Hoàn tất:
- go run ./tools/docs-check, go vet ./..., go test ./... đều pass.
- Commit theo nhóm:
    docs: BD-00, BD-03, DD-00, DD-03 rev 2 (P04b-D1..D6)
    docs: add reconstructed P04 report
    docs: sync README and CLAUDE.md locking scope (P04b-D7, D8)
    chore(docs-check): validate report file names (P04b-D10)
    docs: persist session reports in docs/reports (P04b-D10)
    docs: plan five issues and prompts P06-P10 in P05 (P04b-D11)
    docs: record P04b in progress tracker and prompt index
- Lưu báo cáo của phiên này thành docs/reports/P04b-<YYYY-MM-DD>.md và commit (theo D10), rồi push. Không merge.
- In báo cáo theo mẫu CLAUDE.md §13, tiêu đề "Báo cáo — P04b — <YYYY-MM-DD>". Mục 7 ghi kết quả CI của PR #6.
```
