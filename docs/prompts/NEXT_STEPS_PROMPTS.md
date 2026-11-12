# GitHive — Prompt cho 3 bước tiếp theo

Cách dùng: mỗi bước là **một phiên Claude Code riêng** (`claude` trong thư mục repo). Chỉ bắt đầu bước sau khi PR của bước trước đã được merge vào `main`.

| Bước | Ai làm | Điều kiện bắt đầu | Kết quả |
| --- | --- | --- | --- |
| 1 | Agent (Prompt 1) | Danny đã đặt file thiết kế tổng thể export vào repo (xem bên dưới) | PR `docs: normalize documentation structure` |
| 2a | Agent (Prompt 2) | PR bước 1 đã merge | PR `chore: bootstrap project` |
| 2b | Cố vấn + Danny (song song với 2a) | PR bước 1 đã merge | BD-00, BD-03 → DD-00, DD-03 ở trạng thái "Đã chốt" |
| 3 | Agent (Prompt 3, rồi Prompt mẫu C trong KICKOFF) | PR bước 2a đã merge **và** DD-00, DD-03 đã chốt | Issue + PR cho DD-00 và DD-03 |

**Việc Danny làm trước Prompt 1:** mở tài liệu Basic Design trên claude.ai, chọn **Export → Markdown**, rồi lưu vào repo với tên `docs/basic-design/BD-00-overview.md`. Thư mục `docs/basic-design/` chưa có thì cứ tạo; agent sẽ gộp nó với thư mục cũ. Chép luôn file này vào `docs/Promtps/NEXT_STEPS_PROMPTS.md`; agent sẽ chuyển nó sang `docs/prompts/` cùng các file khác.

---

## Prompt 1 — Chuẩn hóa tài liệu (dán nguyên văn)

```text
Bạn là senior developer của dự án GitHive. Đọc CLAUDE.md trước. Trao đổi với tôi bằng tiếng Việt; mọi thứ ghi vào repo bằng tiếng Anh, TRỪ nội dung requirement gốc, được giữ bằng tiếng Việt như tôi viết.

Phiên này chỉ làm tài liệu, không viết code Go. Tôi (Danny) cho phép trong phiên này: đổi tên và di chuyển file, xóa file .docx, .xlsx và .gitkeep sau khi đã chuyển nội dung sang .md, sửa CLAUDE.md và docs/prompts/KICKOFF_PROMPT.md đúng theo các mục dưới đây. Ngoài các mục đó, quy tắc "Ask Danny first" của CLAUDE.md vẫn áp dụng.

Bước 0 — Kiểm tra:
- Báo đang ở Windows hay WSL. Chạy `git status` và `git log --oneline -5`.
- Kiểm tra file docs/basic-design/BD-00-overview.md đã có chưa. Chưa có thì DỪNG và nhắc tôi export tài liệu Basic Design ra Markdown.
- Tạo nhánh docs/normalize-structure từ origin/main.

Bước 1 — Đổi tên thư mục sang kebab-case, chỉ dùng ký tự ASCII (dùng git mv để giữ lịch sử):
  docs/Requirements   → docs/requirements
  docs/Basic Design   → docs/basic-design   (gộp với file BD-00-overview.md đã có)
  docs/Detail Design  → docs/detail-design
  docs/Promtps        → docs/prompts
Tạo thêm docs/adr/. Xóa các file .gitkeep ở thư mục nào đã có file thật. Thư mục còn rỗng thì giữ .gitkeep.

Bước 2 — Requirement:
- Chuyển docs/requirements/CLI Tool For Worktree.docx thành docs/requirements/README.md. Đọc file bằng công cụ có sẵn (pandoc hoặc python-docx); không có công cụ nào thì báo tôi, đừng tự cài vào hệ thống.
- Giữ nguyên văn tiếng Việt. Chỉ thêm mã và bảng truy vết, KHÔNG sửa nội dung.
- Đánh mã theo đúng bảng sau:
  REQ-01 Chọn/khởi tạo thư mục chứa worktree, có đường dẫn mặc định theo OS  → BD-01
  REQ-02 Chọn thư mục cha chứa các dự án                                     → BD-01
  REQ-03 Chọn dự án làm việc (phải có .git)                                   → BD-02
  REQ-04 Tạo worktree từ commit mới nhất của base branch, nhánh feat/<tên>    → BD-04, BD-03
  REQ-05 Hiển thị mọi worktree của một dự án                                  → BD-05
  REQ-06 Pull base branch rồi rebase/merge vào worktree, cảnh báo conflict    → BD-06
  REQ-07 Xóa nhiều worktree để giải phóng dung lượng                          → BD-07
  REQ-08 Liệt kê dự án và base branch của từng dự án (NULL nếu chưa cấu hình) → BD-02
  NFR-01 Chạy được trên WSL, Ubuntu, Windows, macOS                           → BD-00, BD-09
  NFR-02 Nhiều agent làm song song, mỗi agent một nhánh, không xung đột       → BD-00, BD-04
  NFR-03 Base branch đặt tên khác nhau tùy dự án                              → BD-02
  NFR-04 Cấu hình không mất khi hết phiên hoặc tắt máy                        → BD-01
  NFR-05 Cài đặt và chọn thiết lập qua giao diện tương tác thay vì gõ lệnh    → BD-01, BD-00
- Cuối file thêm bảng "Traceability" với các cột REQ | BD | DD | Issue | Status. Ghi chú riêng rằng BD-08 (path, shell-init) không bắt nguồn từ requirement nào; nó phát sinh từ thiết kế để agent di chuyển vào worktree.
- Sau khi so khớp nội dung Markdown với file .docx (không thiếu đoạn nào), xóa file .docx khỏi repo. Lịch sử git vẫn giữ bản gốc; ghi hash commit chứa bản gốc ở cuối README.md của requirements.

Bước 3 — Danh mục BD:
- Tạo docs/basic-design/README.md gồm:
  - Bảng BD-00…BD-09 với các cột: Mã BD | Chức năng | Lệnh | Mã DD | Trạng thái BD | Trạng thái DD | Phụ thuộc | REQ. Lấy dữ liệu từ file .xlsx; cột REQ lấy từ bảng ở Bước 2. Mọi trạng thái là "Not started".
  - Mục "Conventions": phân biệt BD (làm gì) với DD (làm thế nào), và quy tắc "code chỉ được viết từ DD có trạng thái Approved".
  - Mục "Naming": BD-00-overview.md là bản thiết kế tổng thể, viết khi tool còn tên tạm `wtx`. Đọc wtx = githive, wtx.* = githive.*, WTX_* = GITHIVE_*, new/ls/rm = spawn/swarm/retire. KHÔNG sửa nội dung BD-00-overview.md.
- Mỗi tài liệu BD/DD sau này đặt tên dạng BD-xx-<slug>.md / DD-xx-<slug>.md và mở đầu bằng khối metadata: ID, Title, Status (Draft | In review | Approved), REQ, Depends on, Last updated.
- Xóa file .xlsx khỏi repo. Bảng theo dõi tiến độ nằm trên claude.ai; repo chỉ giữ bản Markdown.

Bước 3b — Thống nhất định dạng: mọi tài liệu trong docs/ là Markdown, theo template.
- Quy tắc: trong docs/ chỉ có file .md. Hình ảnh, nếu có, đặt trong docs/assets/ (.png/.svg); sơ đồ ưu tiên viết bằng Mermaid ngay trong .md. Không còn .docx, .xlsx hay .pdf.
- Tạo docs/templates/ với các file sau, mỗi file mở đầu bằng khối metadata (dùng YAML front matter):
    ---
    id: BD-xx | DD-xx | ADR-xxxx
    title:
    status: Draft            # Draft | In review | Approved | Superseded
    req: [REQ-xx, NFR-xx]
    depends_on: []
    owner:
    last_updated: YYYY-MM-DD
    ---
  1. BD-template.md — các mục: 1 Purpose · 2 Users & use cases · 3 Commands & user-facing I/O (ví dụ lệnh và output) · 4 Main flow (Mermaid flowchart) · 5 Alternate & error flows · 6 Business rules (đánh mã BR-xx) · 7 Platform constraints (Windows / WSL / Linux / macOS) · 8 Out of scope · 9 Open questions · 10 Traceability (REQ → BR).
  2. DD-template.md — các mục: 1 Scope (implements BD-xx, BR nào) · 2 Packages & files · 3 Types & interfaces (Go) · 4 Config keys / JSON output schema · 5 Algorithm: lệnh git chính xác theo thứ tự · 6 Locking & concurrency · 7 Errors & exit codes (bảng) · 8 Platform specifics · 9 Test cases (bảng: TC-xx · Given/When/Then · Type unit/integration/e2e · OS) · 10 Acceptance criteria (checklist, mỗi dòng trỏ tới TC-xx) · 11 Open questions.
  3. ADR-template.md — Context · Decision · Alternatives considered · Consequences · Status.
  4. REQ-format.md — mô tả cách đánh mã REQ-xx/NFR-xx và cách cập nhật bảng Traceability.
  Mỗi mục trong template có một dòng chú thích ngắn (<!-- ... -->) nói phải điền gì. Không điền nội dung mẫu bịa ra.
- BD-00-overview.md: chỉ thêm front matter (id: BD-00-overview, status: Approved, last_updated: 2026-11-12) ở đầu file. Không sửa nội dung; chỗ hình bị thay bằng placeholder khi export thì ghi chú "diagram omitted in export — see claude.ai doc".
- docs/PROGRESS.md, docs/basic-design/README.md, docs/requirements/README.md cũng dùng front matter tương tự (id: PROGRESS, BD-INDEX, REQ).
- Thêm script kiểm tra tools/docs-check (viết bằng Go, chạy bằng `go run ./tools/docs-check`; không dùng bash/PowerShell để chạy được trên mọi OS). Script báo lỗi khi:
    a) trong docs/ có file khác .md, trừ docs/assets/*.png|*.svg;
    b) file BD/DD/ADR thiếu front matter hoặc thiếu trường bắt buộc, hoặc status không thuộc danh sách cho phép;
    c) tên file BD/DD không đúng mẫu ^(BD|DD)-\d{2}-[a-z0-9-]+\.md$ (trừ BD-00-overview.md);
    d) file .md không phải UTF-8 hoặc có CRLF.
  Repo chưa có go.mod thì tạo tools/docs-check/go.mod riêng (module github.com/duong2604/githive/tools/docs-check). Đây là ngoại lệ duy nhất cho quy tắc "không viết code Go" của phiên này; giữ script dưới 200 dòng và có test bảng cho từng quy tắc.
- Chạy script trên repo sau khi chuyển đổi xong; phải sạch lỗi.

Bước 4 — Sửa lỗi mã hóa:
- Ghi lại README.md bằng UTF-8, không BOM, xuống dòng LF. Nội dung tối thiểu: tên dự án, một câu mô tả, trạng thái "pre-alpha", link tới docs/requirements và docs/basic-design.
- Thêm .gitattributes:
    * text=auto eol=lf
    *.docx binary
    *.xlsx binary
    *.png binary
- Chạy `git add --renormalize .` và kiểm tra lại mã hóa của mọi file .md (`file` hoặc PowerShell). Báo file nào không phải UTF-8.

Bước 5 — Cập nhật CLAUDE.md (chỉ sửa đúng các chỗ sau, giữ nguyên phần còn lại):
- §3 Source of truth, thay bằng:
    requirements:    docs/requirements/README.md (REQ-xx, NFR-xx)
    architecture:    docs/basic-design/BD-00-overview.md + docs/basic-design/README.md
    feature designs: docs/basic-design/BD-xx-*.md → docs/detail-design/DD-xx-*.md
    decisions:       docs/adr/
    progress:        docs/PROGRESS.md
  Thêm quy tắc định dạng: docs/ chỉ chứa Markdown theo docs/templates/; tài liệu mới phải tạo từ template; chạy `go run ./tools/docs-check` trước khi mở PR có sửa docs.
  Thêm quy trình: REQ → BD → DD → issue → PR. Chỉ implement từ DD có Status: Approved. Mỗi issue và PR ghi mã DD (ví dụ "Implements DD-04"). DD thiếu hoặc mơ hồ thì hỏi Danny, không tự đoán.
- §7, thêm dòng: commit chỉ có tài liệu dùng `docs:`, cấu hình/công cụ dùng `chore:`; `feat:` chỉ dành cho tính năng người dùng thấy được.
- Thay mọi đường dẫn docs/design/basic-design.md trong file bằng đường dẫn mới.

Bước 6 — Cập nhật docs/prompts/KICKOFF_PROMPT.md:
- Phần B: ghi rõ phiên kickoff chỉ làm Bước 0–2 (dựng khung, CI, testutil, PR).
- Bỏ Bước 3 cũ (tạo issue MVP). Thay bằng: "Issues are created per DD after it is Approved — see docs/prompts/NEXT_STEPS_PROMPTS.md, Prompt 3."
- Sửa đường dẫn tài liệu cho khớp cấu trúc mới.

Bước 7 — Tạo docs/PROGRESS.md (tiếng Anh). Các mục:
- "Current phase": Phase 1 — Documentation baseline.
- "Plan": bảng Step | Owner | Status | Output | Date, gồm các dòng:
    1  Normalize documentation structure | Agent           | In review (khi PR mở) | PR link
    2a Bootstrap project (KICKOFF B, steps 0–2) | Agent    | Not started
    2b BD-00, BD-03 → DD-00, DD-03 approved | Advisor + Danny | Not started
    3  Issues + implementation for DD-00, DD-03 | Agent    | Not started
- "Review findings (2026-11-13)": 7 mục kèm trạng thái Fixed/Open:
    1 missing architecture doc; 2 binary docs; 3 no REQ IDs; 4 folder names; 5 README UTF-16; 6 commit types; 7 kickoff created issues too early; 8 no shared templates or automated docs check.
- "Log": mỗi phiên agent thêm một dòng: ngày, việc đã làm, PR, OS đã kiểm chứng.

Bước 8 — Commit theo từng nhóm logic, ví dụ:
  docs: rename documentation folders to kebab-case
  docs: convert requirements to markdown with REQ IDs
  docs: add basic design index and conventions
  docs: add markdown templates for BD, DD, ADR and requirements
  chore: add docs-check tool
  chore: add .gitattributes and fix README encoding
  docs: update CLAUDE.md and kickoff prompt for BD→DD workflow
  docs: add progress tracker
Push và mở PR (không draft, vì chỉ có tài liệu) với tiêu đề "docs: normalize documentation structure". Body gồm: danh sách thay đổi, bảng ánh xạ các phát hiện review đã sửa, các điểm cần tôi quyết định. Không merge.

Bước 9 — Báo cáo cho tôi bằng tiếng Việt:
- Link PR. Mục nào đã làm, mục nào chưa làm được và vì sao.
- Có đoạn nào trong requirement gốc mơ hồ khi gắn mã REQ không.
- Cập nhật docs/PROGRESS.md trước khi kết thúc.
Dừng ở đây.
```

---

## Prompt 2 — Dựng khung dự án (sau khi PR bước 1 đã merge)

```text
Đọc CLAUDE.md, docs/PROGRESS.md và docs/prompts/KICKOFF_PROMPT.md. Thực hiện đúng phần B, Bước 0–2 của KICKOFF_PROMPT.md. Trong CI, thêm một job chạy `go run ./tools/docs-check` (chỉ cần chạy trên ubuntu). Nếu tạo go.mod ở gốc repo, gộp module tools/docs-check vào module gốc và xóa go.mod riêng của nó. Không tạo issue tính năng và không viết logic tính năng: các phần đó chờ DD được duyệt.
Khi xong: cập nhật docs/PROGRESS.md (bước 2a → In review, thêm link PR, thêm dòng Log), rồi báo cáo bằng tiếng Việt theo Bước 4 của KICKOFF, bỏ phần đề xuất thứ tự issue.
```

---

## Prompt 3 — Tạo issue từ DD đã duyệt (sau khi DD-00, DD-03 đã Approved)

```text
Đọc CLAUDE.md và docs/PROGRESS.md. Với mỗi file trong docs/detail-design/ có Status: Approved mà chưa có issue (xem bảng Traceability trong docs/requirements/README.md):
- Tạo milestone MVP và các label (mvp, v1, v2, design, bug, windows, wsl) nếu chưa có.
- Tạo issue tiêu đề "[DD-xx] <tên>". Body gồm: link tới DD, Acceptance criteria chép đúng từ mục acceptance/test của DD (không tự thêm bớt), danh sách package bị ảnh hưởng, và mục "Can run in parallel with".
- Một DD quá lớn (hơn khoảng 400 dòng code dự kiến) thì đề xuất cách tách thành nhiều issue và hỏi tôi trước khi tạo.
- Mở một PR docs: cập nhật cột Issue trong bảng Traceability và trong docs/basic-design/README.md, cập nhật docs/PROGRESS.md (bước 3 → In progress).
Báo cáo bằng tiếng Việt: danh sách issue, issue nào giao được cho nhiều agent làm song song. Không bắt đầu code.
```

Sau đó, mỗi issue được giao bằng **Prompt mẫu C** trong `docs/prompts/KICKOFF_PROMPT.md`. Bổ sung vào đầu prompt dòng: "Implements DD-xx — đọc file DD trước, chỉ làm đúng phạm vi DD."
