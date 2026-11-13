---
id: PO-01
title: PO / TL / Dev operating model and DD rev 3
from: PO
to: TL
step: "3"
status: Done
branch: chore/po-01-operating-model
pr: "#7"
depends_on: [P04b]
decisions: [PO-01-D1, PO-01-D2, PO-01-D3, PO-01-D4, PO-01-D5, PO-01-D6]
last_updated: 2026-11-13
---

# PO-01 — Mô hình làm việc PO / TL / Dev và DD rev 3

Nguyên văn prompt Danny dán vào phiên (R3). Các ký tự lạ như `</nn>`, `</slug>` có sẵn trong bản dán và được giữ nguyên.

## Prompt (nguyên văn)

```text
PO-01 — Mô hình làm việc PO / TL / Dev và DD rev 3
Người gửi: PO (Claude trên claude.ai) · Người nhận: Technical Leader (TL) · Danny dán nguyên văn toàn bộ khối này.

==================================================================
VAI TRÒ CỦA BẠN
==================

Bạn là Technical Leader (TL) của dự án GitHive (github.com/duong2604/githive). Đội gồm:

- Danny (duong2604): chủ dự án. Trả lời câu hỏi về requirement, duyệt và merge PR. Là người chuyển prompt giữa các bên.
- PO (Claude trên claude.ai): sở hữu requirement, đặc tả sản phẩm và UI, tiêu chí nghiệm thu. Giao việc cho bạn bằng các prompt có mã PO-xx.
- TL (bạn):
  - Sở hữu thiết kế kỹ thuật (DD, ADR) và các quyết định kỹ thuật.
  - Chia việc, viết prompt cho Dev (P05, P06…), review mọi PR của Dev trước khi Danny merge.
  - Báo cáo tiến độ cho PO.
- Dev agent: viết code theo prompt của bạn. Báo cáo cho bạn.
  Trước đây "Advisor" (cũng là Claude trên claude.ai) kiêm cả hai việc sản phẩm và kỹ thuật. Từ PO-01 trở đi: phần sản phẩm thuộc PO, phần kỹ thuật thuộc bạn. Các quyết định cũ ghi "Advisor" vẫn giữ nguyên giá trị.

==================================================================
QUY TẮC GIỮ LỊCH SỬ (thay cho việc chép file thủ công)
==============================================================

R1. Không ai được tải file về rồi chép đè vào repo nữa. Mọi thay đổi đều do agent làm trên một nhánh, commit và mở PR.
R2. Nội dung do PO soạn chỉ đến qua prompt. Khi tài liệu đã có thì PO gửi yêu cầu sửa theo từng mục; chỉ khi tạo tài liệu mới mới gửi toàn văn.
R3. Commit đầu tiên của mỗi phiên PO-xx: lưu nguyên văn prompt vào docs/prompts/PO-<NN></nn>-<slug></slug>.md. Lấy text Danny dán; thêm front matter như prompt P0x, với id: PO-<NN></nn>, from: PO, to: TL.
R4. Prompt bạn viết cho Dev được tạo trong repo (docs/prompts/P<NN></nn>-<slug></slug>.md, from: TL, to: Dev) và commit qua PR. Danny chỉ việc mở file và dán; không phải chép file.
R5. Mỗi phiên lưu báo cáo vào docs/reports/<ID></id>-<YYYY-MM-DD></yyyy>.md (P04b-D10). Áp dụng cho cả mã PO-xx.
R6. Báo cáo theo mốc công việc:

- Kết thúc mỗi prompt hoặc phiên: báo cáo theo R5.
- Kết thúc mỗi issue: Dev ghi báo cáo issue vào chính PR của issue đó (docs/reports/P<NN></nn>-<ngày>.md).
- Khi mọi phần của một DD đã merge, nghĩa là tính năng xong: TL viết báo cáo tính năng docs/reports/FEAT-DD-<NN></nn>-<YYYY-MM-DD></yyyy>.md, gồm: tính năng làm được gì so với REQ/BD, đối chiếu từng tiêu chí nghiệm thu, các điểm lệch, và việc còn lại.
  R7. File bàn giao luôn sẵn sàng, để một agent khác thay thế ngay khi agent hiện tại offline:
- Mỗi agent có một file bàn giao trên nhánh đang làm: docs/handoff/TL.md, docs/handoff/DEV-1.md, DEV-2.md… Mỗi Dev chạy song song có một mã riêng; Danny ghi mã này ở dòng đầu prompt.
- File theo mẫu docs/templates/HANDOFF-template.md, gồm:
  - vai trò và mã agent, prompt hoặc issue đang làm, nhánh, PR, commit cuối;
  - mục tiêu của phiên;
  - đã xong;
  - đang dở, chính xác tới file và hàm;
  - bước tiếp theo, đánh số;
  - lệnh để tiếp tục (git switch …, go test …);
  - câu hỏi đang chờ, gửi ai;
  - các quyết định đã đưa ra trong phiên;
  - thời điểm cập nhật cuối.
- Cập nhật và push file bàn giao: sau mỗi nhóm commit, trước mỗi lần chạy lệnh dài hoặc chờ CI, và trước khi kết thúc phiên. Không bao giờ để quá một nhóm commit mà file bàn giao chưa được cập nhật.
- Ngay sau commit đầu tiên của phiên, mở PR ở dạng draft, để `gh pr list` luôn cho thấy ai đang làm gì, trên nhánh nào.
- Khi PR merge, file bàn giao vào main cùng PR, và trở thành dấu vết lần cuối của agent đó.
  R8. Quy trình thay thế agent:
- Agent mới chạy: gh pr list --state open; git switch <nhánh>; git pull; đọc docs/handoff/<mã>.md, CLAUDE.md và báo cáo gần nhất trong docs/reports/.
- Sau đó ghi một dòng "Taken over by <mã mới> at <thời gian>" vào file bàn giao, rồi làm tiếp từ "Bước tiếp theo".
- Không làm lại phần đã commit.

==================================================================
BƯỚC 0 — Dọn thay đổi thủ công còn treo
================================================

Trên main ở máy này đang có 3 file Danny chép tay, chưa commit:
docs/detail-design/DD-00-foundation.md, docs/detail-design/DD-03-target.md, docs/prompts/P05-create-issues-from-dd.md.

1. git switch main && git pull --ff-only. Đảm bảo main có ab39994 (PR #6 đã merge).
2. git status --short và git diff --stat: ghi lại kết quả vào báo cáo.
3. Bỏ 3 thay đổi đó: git restore docs/detail-design/DD-00-foundation.md docs/detail-design/DD-03-target.md docs/prompts/P05-create-issues-from-dd.md (Danny cho phép). Nội dung tương ứng được áp dụng lại ở Bước 3 dưới dạng sửa từng mục.
4. git status --short phải rỗng. Nếu còn file nào khác thì dừng lại và báo cáo.

==================================================================
BƯỚC 1 — Nhánh và lưu prompt
==================================

Tạo nhánh chore/po-01-operating-model từ origin/main. Commit file prompt theo R3 với tên docs/prompts/PO-01-operating-model.md.

==================================================================
BƯỚC 2 — Ghi mô hình làm việc vào repo
==============================================

2a. CLAUDE.md §1 "Team and roles": thay bảng bằng 4 vai trò Danny / PO / TL / Dev theo mô tả ở trên. Thêm một câu: "Advisor decisions recorded before PO-01 remain valid."
2b. CLAUDE.md thêm §14 "Working agreement", gồm:

- Bảng sở hữu tài liệu:
  - docs/requirements và docs/product: PO
  - docs/basic-design (hành vi người dùng thấy): PO viết, TL review tính khả thi
  - docs/detail-design và docs/adr: TL
  - code và test: Dev
  - review PR: TL
  - merge: Danny
- Luồng làm việc: PO-xx → TL → (P-xx) → Dev → báo cáo → TL review → báo cáo TL kèm "PO summary" → Danny chuyển cho PO.
- Quy tắc R1–R5 ở trên.
- Cách xếp câu hỏi: câu hỏi kỹ thuật thì TL tự quyết; câu hỏi sản phẩm hoặc UX thì gửi PO; câu hỏi về business hoặc requirement gốc thì gửi Danny.
  2c. CLAUDE.md §13: bổ sung rằng báo cáo của TL gửi PO có thêm mục 9 "PO summary": tối đa 10 dòng, viết bằng ngôn ngữ nghiệp vụ, gồm: tính năng nào xong, đúng hay lệch so với requirement ở chỗ nào, cần PO hoặc Danny quyết định gì.
  2c-bis. CLAUDE.md thêm §15 "Handoff and takeover", ghi nguyên R6–R8 (bằng tiếng Anh). Tạo docs/templates/HANDOFF-template.md có đủ các mục ở R7. Tạo docs/handoff/TL.md cho chính phiên này và cập nhật theo R7. Đây là file bàn giao đầu tiên.
  2d. docs/prompts/README.md:
- Thêm cột "From → To".
- Thêm quy ước tên PO-<NN></nn>[a-z]-<slug></slug>.md.
- Thêm dòng PO-01.
  2e. tools/docs-check:
- Chấp nhận tên prompt ^(P\d{2}|PO-\d{2})[a-z]?-[a-z0-9-]+\.md$.
- Chấp nhận tên báo cáo ^(P\d{2}|PO-\d{2})[a-z]?-\d{4}-\d{2}-\d{2}(-\d+)?\.md$.
- Chấp nhận thêm tên báo cáo tính năng ^FEAT-DD-\d{2}-\d{4}-\d{2}-\d{2}(-\d+)?\.md$.
- Chấp nhận thư mục docs/product/ (chỉ áp quy tắc chung: .md, UTF-8, LF).
- Thư mục docs/handoff/: tên file phải khớp ^(TL|DEV-\d+)\.md$.
- Thêm test bảng cho các quy tắc này.

==================================================================
BƯỚC 3 — DD-00 và DD-03 rev 3 (sửa từng mục, không chép đè)
======================================================================

Các quyết định dưới đây do Advisor đưa ra khi review báo cáo P04b. Lưu với mã PO-01-D1…D6. Bạn là TL nên được phép phản biện: nếu không đồng ý điểm nào, vẫn áp dụng, rồi ghi phản biện vào báo cáo.
PO-01-D1  DD-00 §1, bảng "Implementation slices":

- Thêm cột "Test cases": 00a = TC-01…TC-11, TC-25…TC-27; 00b = TC-12…TC-18; 00c = TC-19…TC-24, TC-28.
- Ô Size của 00a ghi: "~450 lines incl. tests (accepted above the ~400 guide, PO-01-D1)".
  PO-01-D2  Issue lấy tiêu chí nghiệm thu theo cột "Test cases" đó.
  PO-01-D3  DD-03 §7, bảng lỗi: thêm dòng ngay sau dòng "`project@*` where not allowed":
  | A single task given to `swarm` (`AllowTask` false) | `invalid_target` | `<cmd> takes 'project@*' or --all, not a single task` | 2 |
  PO-01-D4  DD-00 §3.5: thay comment trên `func With(` thành:
  // With acquires, runs fn, and always releases. Errors are returned through apperr.From,
  // so a cancelled ctx (in Acquire or from fn) becomes interrupted (BR-12).
  PO-01-D5  DD-00 §9, TC-26 sửa thành:
  `NoPrompter(false)`: `Input`, `Select` and `Confirm` → `input_required`; `NoPrompter(true)`: `Confirm` → true, `Input` and `Select` still → `input_required` (BR-02)
  Câu `repo_on_mnt` trong DD-03 §5.4 là câu chuẩn; ví dụ trong BD-00 chỉ để minh họa.
  PO-01-D6  Báo cáo đã lưu trong repo không cần sửa lại để thêm kết quả CI, vì kết quả CI xem trên PR.
  Mỗi file thêm một dòng vào "Revision history":
  "| 2026-11-13 | Rev 3 (PO-01): <tóm tắt> |"

==================================================================
BƯỚC 4 — Tiếp nhận vai trò TL (chỉ đọc, đánh giá, không code)
==========================================================================

Đọc: CLAUDE.md, docs/requirements/README.md, docs/basic-design/*, docs/detail-design/*, docs/adr/*, docs/PROGRESS.md, docs/prompts/README.md, docs/prompts/P05-create-issues-from-dd.md, docs/reports/*.
Trong báo cáo, viết mục "Đánh giá của TL":
  a) Rủi ro kỹ thuật lớn nhất trước khi code, tối đa 5 điểm.
  b) P05 (tạo 5 issue) có phù hợp mô hình mới không. Từ nay P05 là prompt bạn tự chạy hoặc giao cho Dev. Đề xuất sửa P05 nếu cần, nhưng chưa sửa và chưa chạy.
  c) Kế hoạch giao việc cho Dev (P06…P10): thứ tự, cái gì làm song song, bạn sẽ review gì ở mỗi PR.
Không tạo issue, không viết code ở phiên này.

==================================================================
BƯỚC 5 — Hoàn tất
======================

- go run ./tools/docs-check, go vet ./..., go test ./... đều pass.
- Commit theo nhóm:
  docs: save PO-01 prompt
  docs: PO/TL/Dev roles and working agreement in CLAUDE.md
  docs: handoff and takeover rules, handoff template (R6-R8)
  chore(docs-check): accept PO prompts, PO reports and docs/product
  docs: DD-00, DD-03 rev 3 (PO-01-D1..D6)
  docs: record PO-01 in progress tracker and prompt index
- docs/PROGRESS.md: thêm bước "PO-01" (In review); thêm dòng Log ghi rõ đã bỏ 3 file chép tay ở Bước 0.
- Mở PR ở dạng draft ngay sau commit đầu tiên (R7), tiêu đề "chore: PO/TL/Dev operating model and DD rev 3 (PO-01)". Khi xong và CI xanh, chuyển sang ready. Không merge.
- Cập nhật docs/handoff/TL.md lần cuối, ghi trạng thái "waiting for Danny to merge".
- Lưu báo cáo vào docs/reports/PO-01-<YYYY-MM-DD></yyyy>.md (commit, push). In báo cáo theo mẫu §13, tiêu đề "Báo cáo — PO-01 — <ngày>". Báo cáo có thêm mục "Đánh giá của TL" (Bước 4) và mục 9 "PO summary".
```
