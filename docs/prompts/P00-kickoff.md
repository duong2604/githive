---
id: P00
title: Kickoff (machine setup, bootstrap, task template)
step: "0"
status: Ready
branch: ""
pr: ""
depends_on: []
decisions: []
last_updated: 2026-11-13
---

# Prompt cho Claude Code — GitHive

Gồm 3 phần:
- **A.** Chuẩn bị máy (Danny làm, một lần cho mỗi máy).
- **B.** Prompt phiên kickoff: chỉ làm Bước 0–2 (dựng khung, CI, testutil, PR). Không tạo issue.
- **C.** Prompt mẫu cho mỗi task sau này.

---

## A. Chuẩn bị máy

**Cả hai máy:** cài Go (bản stable mới nhất), Git ≥ 2.38, GitHub CLI (`gh auth login`), Claude Code, golangci-lint.

**Máy cá nhân (Windows):**
```powershell
git config --global core.longpaths true
cd D:\workspace          # hoặc %USERPROFILE%\workspace nếu không có ổ D:
git clone https://github.com/duong2604/githive.git
cd githive
```

**Laptop công ty (WSL):**
```bash
mkdir -p ~/workspace && cd ~/workspace      # KHÔNG dùng /mnt/c hay /mnt/d
git clone https://github.com/duong2604/githive.git
cd githive
```

Sau đó:
1. Chép `CLAUDE.md` vào gốc repo. Phiên đầu tiên sẽ commit file này trong PR bootstrap.
2. Mở tài liệu Basic Design, chọn **Export → Markdown**, lưu thành `docs/basic-design/BD-00-overview.md` (đã làm trong PR `docs: normalize documentation structure`).
3. Chạy `claude` trong thư mục repo rồi dán prompt ở phần B.

---

## B. Prompt phiên đầu tiên (dán nguyên văn)

```text
Bạn là senior developer của dự án GitHive. Tôi (Danny) là khách hàng/product owner; cố vấn kỹ thuật là người viết docs/basic-design/BD-00-overview.md. Đọc kỹ CLAUDE.md và docs/PROGRESS.md trước khi làm bất cứ điều gì và tuân thủ tuyệt đối. Trao đổi với tôi bằng tiếng Việt; mọi thứ ghi vào repo bằng tiếng Anh.

Repo hiện chỉ có tài liệu trên main (CLAUDE.md, README.md, docs/, tools/docs-check). Nhiệm vụ phiên này: chỉ làm Bước 0–2 (dựng khung, CI, testutil, PR). KHÔNG tạo issue và CHƯA implement tính năng.

Bước 0 — Kiểm tra môi trường, báo cáo rồi mới làm tiếp:
- Đang ở Windows hay WSL? Đường dẫn repo có nằm dưới /mnt/* không?
- Phiên bản go, git (cần ≥ 2.38), gh (đã auth chưa), golangci-lint. Thiếu gì thì liệt kê lệnh cài, đừng tự cài vào hệ thống.
- Hỏi tôi đúng 1 câu trước khi tạo file: chọn license MIT hay Apache-2.0.

Bước 1 — Nhánh chore/bootstrap tạo từ origin/main. Tạo:
- go.mod (module github.com/duong2604/githive), cmd/githive/main.go với cobra root, lệnh `version` (in version/commit/date qua ldflags) và `--json`.
- Bộ khung thư mục theo CLAUDE.md §5. Chỉ tạo package khi có code thật; không tạo file rỗng cho có.
- .gitattributes (`* text=auto eol=lf`), .gitignore, .editorconfig, .golangci.yml (bộ linter vừa phải).
- .github/workflows/ci.yml: matrix ubuntu-latest, windows-latest, macos-latest; các bước vet, golangci-lint, test (thêm -race trên Linux/macOS), build. Cài git đủ mới trên runner nếu cần.
- .goreleaser.yaml: build windows/linux/darwin × amd64/arm64, archive kèm binary git-hive. Chỉ cấu hình, KHÔNG publish.
- internal/testutil: hàm dựng bare remote có nhánh main và dev + một clone, cô lập hoàn toàn khỏi git config của máy (CLAUDE.md §10), kèm một smoke test dùng nó.
- README.md tiếng Anh: tầm nhìn, trạng thái "pre-alpha", bảng lệnh dự kiến (CLAUDE.md §4), cách build.
- CONTRIBUTING.md ngắn, CHANGELOG.md (mục Unreleased), LICENSE theo lựa chọn của tôi.
- .github/pull_request_template.md, .github/ISSUE_TEMPLATE/{feature,bug}.md.
- docs/adr/0001-record-architecture-decisions.md và 0002-core-technical-decisions.md (Approved), tóm tắt CLAUDE.md §4–5: Go, git CLI, config 2 tầng, khóa theo repo, cú pháp project@task, tên lệnh spawn/swarm/retire kèm alias.
- CLAUDE.md và docs/ đã có trên main; chỉ cập nhật README.md, và docs/PROGRESS.md theo docs/prompts/P02-bootstrap-project.md.

Bước 2 — Mở draft PR "chore: bootstrap project", chờ CI xanh trên cả 3 OS. CI đỏ thì tự sửa. Xanh rồi thì chuyển sang ready for review. Không merge.

Bước 3 — Issues are created per DD after it is Approved — see docs/prompts/P03-create-issues-from-dd.md.

Bước 4 — Báo cáo cho tôi bằng tiếng Việt:
- Link PR.
- Đã kiểm chứng gì trên OS nào.
- Rủi ro hoặc điểm chưa chắc chắn.

Dừng sau bước 4.
```

---

## C. Prompt mẫu cho mỗi task

```text
Làm issue #<N> của GitHive. Implements DD-xx — đọc CLAUDE.md, file DD và issue (gh issue view <N>) trước; chỉ làm đúng phạm vi DD.
- Làm trong worktree riêng, nhánh feat/<N>-<slug> tạo từ origin/main. Không đụng worktree hoặc nhánh khác.
- Lên kế hoạch ngắn (các file sẽ sửa, test sẽ viết) rồi làm luôn. Chỉ dừng hỏi khi gặp mục "Ask Danny first" hoặc cần ADR.
- Viết test trước với parser/logic thuần; integration test dùng internal/testutil.
- Chạy vet, lint, test. Mở draft PR, chờ CI xanh trên 3 OS, rồi chuyển sang ready.
- Báo cáo bằng tiếng Việt theo Definition of Done trong CLAUDE.md §11.
```

Khi chạy nhiều agent song song, giao mỗi agent một issue ít đụng nhau, theo mục "Can run in parallel with" trong issue (ví dụ issue của DD-03 target và DD-00 git runner). Mỗi agent mở một terminal riêng, chạy `claude --worktree issue-<N>` rồi dán prompt C.
