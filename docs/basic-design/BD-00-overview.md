# wtx — Basic Design: CLI quản lý Git Worktree cho nhiều Claude agent

Nov 12, 2026 · @danny

## 1. Đánh giá requirement

Hướng đi đúng: `git worktree` chính là lời giải cho bài toán "chỉ một agent được checkout". Một repo, nhiều thư mục làm việc, mỗi thư mục một nhánh, dùng chung object database. Tool của bạn có giá trị ở chỗ **chuẩn hóa quy trình** (luôn tạo từ commit mới nhất của nhánh base, đặt tên thống nhất, dọn dẹp an toàn), không phải ở bản thân lệnh git.

### Nguyên nhân gốc của vấn đề hiện tại

Một working directory chỉ có một `HEAD` và một index. Bốn agent dùng chung thư mục nên `git checkout` của agent này phá file đang sửa của agent khác. Worktree tách riêng `HEAD`, index và file, nên mỗi agent có sandbox riêng. Git còn tự cấm một nhánh được checkout ở hai worktree, tức là quy tắc "mỗi agent một nhánh" được git cưỡng chế sẵn.

### Những điểm requirement đúng và nên giữ

- Base branch cấu hình theo từng dự án (`dev`, `developer`, `develop`…).
- Tên nhánh = tiền tố + tên worktree (`feat/worktree-18`).
- Có lệnh list, sync, xóa hàng loạt, list dự án.
- Cấu hình phải bền vững qua phiên làm việc và khởi động lại máy.

### Những điểm còn thiếu hoặc có rủi ro

| # | Vấn đề | Vì sao quan trọng | Đề xuất |
| --- | --- | --- | --- |
| 1 | Đường dẫn `D:\workspace\worktrees\worktree-18` không có tên dự án | Hai dự án cùng có `worktree-18` sẽ đụng nhau | `<worktrees_root>/<project>/<name>` |
| 2 | "Lấy commit mới nhất từ dev" chưa nói rõ local hay remote | `dev` local thường cũ; không cần (và không nên) checkout `dev` | `git fetch origin dev` rồi tạo từ `origin/dev` |
| 3 | Git vẫn có tài nguyên dùng chung giữa các worktree | Hai agent cùng `fetch`/`worktree add` có thể lỗi `cannot lock ref` | Khóa file theo repo, chỉ giữ trong vài giây |
| 4 | WSL vs Windows | Repo nằm ở `/mnt/d` chạy từ WSL rất chậm; worktree tạo bằng git WSL ghi đường dẫn Linux, git Windows không đọc được | Chọn một phía cho mỗi repo; `doctor` cảnh báo |
| 5 | File không được track (`.env`, `appsettings.Development.json`) | Worktree mới không có chúng, agent build lỗi | Cấu hình `copy_files` + hook `post_create` |
| 6 | Agent cũng là người dùng tool | Wizard tương tác làm agent bị treo | Mọi lệnh chạy được bằng flag, có `--json`, exit code rõ ràng |
| 7 | Xóa worktree an toàn | GitLab squash merge làm `git branch --merged` không nhận ra nhánh đã merge | Kiểm tra dirty, commit chưa push, nhánh remote đã bị xóa |
| 8 | Windows MAX\_PATH 260 ký tự | `node_modules` sâu trong worktree sẽ lỗi | Đường dẫn gốc ngắn, bật `core.longpaths` |
| 9 | Cổng dev server | 4 agent cùng chạy `npm run dev` sẽ tranh port 3000 | Để ngoài MVP; cân nhắc cấp port theo worktree ở v2 |

## 2. Định vị so với `claude --worktree`

Claude Code đã có sẵn cờ `--worktree` (`-w`) từ bản 2.1.49 (tháng 2/2026), tạo worktree trong `.claude/worktrees/` của repo và có file `.worktreeinclude` để chép file bị gitignore ([docs](https://code.claude.com/docs/en/worktrees)). Vì vậy `wtx` không nên cạnh tranh việc "tạo worktree cho agent", mà làm phần Claude Code không làm:

- Fetch và tạo từ `origin/<base>` cấu hình theo dự án, đặt tên nhánh theo quy ước team (`feat/...`).
- Quản lý nhiều dự án trong một workspace, xem tổng thể mọi worktree cùng trạng thái ahead/behind/dirty.
- Sync với base có dự đoán conflict trước, dọn dẹp an toàn theo lô.
- Dùng được cho cả agent khác và người làm tay.

Mô hình khuyến nghị: `wtx new task-18` tạo worktree, sau đó `cd` vào và chạy `claude` bình thường (không dùng `-w`). V2 có thể thêm `wtx new task-18 --run claude` để mở agent luôn. Cần kiểm tra trước khi code: `-w` tạo từ nhánh nào và có cấu hình được base không; nếu đã đủ cho nhu cầu của bạn thì phạm vi `wtx` thu hẹp còn sync + cleanup + multi-project.

## 3. Quyết định kiến trúc chính

| Quyết định | Lựa chọn | Lý do | Phương án bị loại |
| --- | --- | --- | --- |
| Ngôn ngữ | **Go** (cobra + charmbracelet/huh cho prompt) | Một file binary, cross-compile sẵn Windows/Linux/macOS × amd64/arm64, khởi động nhanh, không cần runtime | Node/TS (phải cài Node), Rust (đường cong học dốc hơn), bash/PowerShell (không chạy chung được) |
| Tương tác git | Gọi **git CLI** có sẵn trên máy, parse output `--porcelain -z` | Dùng đúng credential GitLab, hook, config của user; worktree support đầy đủ | go-git (hỗ trợ worktree hạn chế, không dùng credential helper) |
| Nguồn sự thật về worktree | **`git worktree list`**, không tự lưu danh sách | Không bao giờ lệch với thực tế khi user/agent dùng git tay | File state/SQLite riêng (dễ drift) |
| Cấu hình toàn cục | File TOML trong thư mục config chuẩn của OS | Bền qua tắt máy, user sửa tay được | Biến môi trường (mất khi đóng shell) |
| Cấu hình theo dự án | **`git config --local wtx.*`** trong `.git/config` của repo | Đi theo repo, xóa repo là xóa config, mọi worktree đọc chung | Map project → config trong file global (lệch khi đổi tên/di chuyển repo) |
| Đồng thời | Khóa file `wtx.lock` trong `git rev-parse --git-common-dir`, timeout 30 giây | Chặn 2 agent fetch/add/remove cùng lúc trên một repo; repo khác không chờ nhau | Không khóa (lỗi ngẫu nhiên `cannot lock ref`) |
| Chế độ chạy | Thiếu tham số + có TTY → hỏi tương tác; không TTY hoặc `--yes` → báo lỗi rõ ràng | Người dùng thích menu, agent cần lệnh xác định | Chỉ TUI |
| Phân phối | GoReleaser → GitHub/GitLab Release, Homebrew tap, Scoop/winget, script `install.sh` | Cài 1 lệnh trên mọi OS | — |

Phiên bản git tối thiểu đề xuất: 2.38 (có `git merge-tree --write-tree` để dự đoán conflict). Git 2.48+ có thêm `worktree.useRelativePaths`, hữu ích khi di chuyển thư mục hoặc dùng lẫn WSL/Windows.

## 4. Quy tắc chọn đường dẫn

Trả lời câu hỏi "ổ D: hay C:": **đừng chọn ổ đĩa, hãy buộc worktree nằm cùng chỗ với thư mục dự án.** User chọn projects root một lần; worktrees root mặc định suy ra từ đó. Như vậy câu hỏi D/C tự biến mất, và worktree luôn cùng ổ/cùng filesystem với repo (nhanh, không lỗi hardlink, không lẫn Windows/WSL).

### Thứ tự ưu tiên khi xác định đường dẫn

1. Flag trên lệnh (`--path`).
2. Biến môi trường (`WTX_WORKTREES_ROOT`).
3. Cấu hình theo repo (`wtx.worktreesRoot`).
4. Cấu hình toàn cục (`worktrees_root`).
5. Mặc định tính toán: `<projects_root>/worktrees`.

### Đề xuất mặc định trong wizard `wtx init` (luôn hiển thị để user xác nhận hoặc sửa)

| Môi trường | Projects root gợi ý | Ghi chú |
| --- | --- | --- |
| Windows có ổ D: cố định, ghi được | `D:\workspace` | Kiểm tra loại ổ là fixed (không phải USB/ổ mạng/DVD) |
| Windows chỉ có C: | `%USERPROFILE%\workspace` | Không dùng `C:\workspace` gốc ổ: máy công ty có thể chặn quyền tạo thư mục gốc |
| WSL | `~/workspace` (filesystem ext4 của WSL) | **Không** gợi ý `/mnt/d/...`: I/O qua 9P chậm hơn nhiều lần, npm/dotnet build rất chậm |
| Linux / macOS | `~/workspace` | — |

Cấu trúc thư mục kết quả (ví dụ Windows):

```
D:\workspace\
  kirindo\            ← repo chính (có .git)
  crm2\
  worktrees\
    kirindo\
      worktree-18\    ← nhánh feat/worktree-18
      worktree-19\
    crm2\
      fix-login\
```

### Quy tắc phụ

- Không đặt worktree bên trong repo chính (IDE, `dotnet build`, `tsc` quét cả thư mục con).
- `worktrees` bị bỏ qua khi quét dự án (không tính là một project).
- Windows: đường dẫn gốc càng ngắn càng tốt; `doctor` kiểm tra `core.longpaths=true`.
- WSL: nếu repo nằm ở `/mnt/*`, `doctor` cảnh báo. Mỗi repo chỉ thao tác bằng một git (WSL hoặc Windows), vì file `.git` trong worktree chứa đường dẫn tuyệt đối của phía đã tạo nó.
- Tên worktree: chỉ `[a-z0-9._-]`, tối đa 50 ký tự, kiểm tra thêm bằng `git check-ref-format --branch`.

## 5. Kiến trúc tầng & module

&#91;embedded content: Kiến trúc wtx · 4 tầng, 3 tài nguyên ngoài\]

CLI chỉ phân tích tham số và gọi service; service chứa logic nghiệp vụ và không biết đến terminal; tầng hạ tầng là nơi duy nhất chạm vào git, file hệ thống và OS. Nhờ vậy test service bằng repo git tạm trong thư mục temp, chạy giống nhau trên CI Windows/Linux/macOS.

Cấu trúc thư mục đề xuất:

```
wtx/
  cmd/wtx/main.go
  internal/cli/        init.go, project.go, new.go, ls.go, sync.go, rm.go, doctor.go
  internal/ui/         prompt.go (huh), table.go, json.go
  internal/service/    project.go, worktree.go, sync.go, cleanup.go
  internal/git/        runner.go, worktree_parse.go, status.go
  internal/config/     global.go (TOML), repo.go (git config wtx.*), resolve.go
  internal/lock/       filelock.go (flock / LockFileEx)
  internal/platform/   paths.go, wsl.go, drives_windows.go
  .goreleaser.yaml
```

## 6. Mô hình cấu hình & dữ liệu

Hai tầng cấu hình, không có database riêng. Trạng thái worktree luôn đọc trực tiếp từ git, nên xóa/tạo worktree bằng git tay cũng không làm tool sai.

### Tầng 1 — Cấu hình toàn cục (bền qua tắt máy)

| OS | Vị trí file |
| --- | --- |
| Windows | `%APPDATA%\wtx\config.toml` |
| Linux / WSL | `$XDG_CONFIG_HOME/wtx/config.toml` (mặc định `~/.config/wtx/config.toml`) |
| macOS | `~/.config/wtx/config.toml` (quen thuộc với dev hơn `~/Library/Application Support`) |

```toml
version = 1
projects_root  = "D:/workspace"
worktrees_root = "D:/workspace/worktrees"   # bỏ trống = <projects_root>/worktrees

[defaults]
remote        = "origin"
branch_prefix = "feat/"
sync_strategy = "rebase"     # rebase | merge
copy_files    = [".env", ".env.local"]
```

Ghi file kiểu atomic (ghi file tạm rồi rename) để không hỏng config khi mất điện giữa chừng. Có trường `version` để migrate về sau.

### Tầng 2 — Cấu hình theo dự án (trong `.git/config` của repo)

```ini
[wtx]
    baseBranch   = dev          ; chưa set → hiển thị NULL
    branchPrefix = feat/        ; ghi đè defaults
    syncStrategy = merge
    copyFiles    = .env
    copyFiles    = appsettings.Development.json
    postCreate   = npm ci
```

Set bằng `wtx project config kirindo --base dev` (bên dưới gọi `git config --local wtx.baseBranch dev`). Vì nằm trong `.git/config` chung, mọi worktree của repo đọc cùng giá trị. Team muốn chia sẻ quy ước thì có thể commit thêm file `.wtx.toml` ở gốc repo (đọc với ưu tiên thấp hơn `.git/config`) — để v2.

### Thứ tự merge cấu hình

Flag → biến môi trường `WTX_*` → `.git/config` (`wtx.*`) → `.wtx.toml` trong repo (v2) → `config.toml` toàn cục → giá trị mặc định của tool.

### Thực thể runtime (không lưu, tính mỗi lần chạy)

| Thực thể | Trường chính | Nguồn |
| --- | --- | --- |
| Project | name, path, commonDir, remote, baseBranch (nullable), worktreeCount | Quét `projects_root` độ sâu 1 + `git config` |
| Worktree | name, path, branch, head, isMain, locked, prunable | `git worktree list --porcelain -z` |
| WorktreeStatus | dirty, ahead, behind (so với `origin/<base>`), upstreamGone, lastCommitAt, diskSize | `git status --porcelain`, `rev-list --left-right --count`, `for-each-ref` |

## 7. Thiết kế lệnh CLI

Tên `wtx` là tạm. Mọi lệnh nhận `-p/--project` (bỏ qua nếu đang đứng trong repo hoặc worktree của nó), `--json`, `--yes`. Chạy `wtx` không tham số trong terminal = menu tương tác.

| Lệnh | Requirement | Mô tả |
| --- | --- | --- |
| `wtx init` | 4.1, 4.2 | Wizard: projects root, worktrees root, prefix, chiến lược sync. Chạy lại để sửa |
| `wtx project list` | 4.8 | Bảng: tên dự án, base branch (hoặc `NULL`), số worktree, remote |
| `wtx project config <p> --base dev [--prefix feat/] [--sync merge]` | 4.3, 4.8 | Set cấu hình theo dự án; không truyền `--base` thì gợi ý từ `origin/HEAD` và các nhánh `dev`/`develop`/`developer` có trên remote |
| `wtx new <name> [--base x] [--branch y] [--no-hooks]` | 4.4 | Fetch base, tạo worktree + nhánh `feat/<name>`, chép file, chạy hook, in đường dẫn |
| `wtx new <name> --from-remote <branch>` | — | Mở worktree cho nhánh đã có trên remote (review MR của người khác) |
| `wtx ls [--all]` | 4.5 | Worktree của dự án (hoặc tất cả): tên, nhánh, ahead/behind, dirty, lần commit cuối |
| `wtx sync [<name>] [--rebase\|--merge] [--dry-run] [--autostash]` | 4.6 | Fetch base, dự đoán conflict, rồi rebase/merge |
| `wtx rm <names...>` / `wtx rm --merged` / `wtx rm --stale 14d` / `wtx rm -i` | 4.7 | Xóa worktree + nhánh local, có kiểm tra an toàn, `--force` để bỏ qua |
| `wtx path <name>` | — | In đường dẫn, dùng `cd $(wtx path x)`; kèm `wtx shell-init` để có `wcd x` |
| `wtx doctor` | — | Kiểm tra git version, `core.longpaths`, repo trên `/mnt/*`, worktree mồ côi (`git worktree prune`) |
| `wtx config get/set <key>` | — | Xem/sửa config toàn cục không cần mở file |

### Exit code (để agent phân nhánh xử lý)

| Code | Ý nghĩa |
| --- | --- |
| 0 | Thành công |
| 1 | Lỗi chung |
| 2 | Sai tham số / thiếu cấu hình (ví dụ base branch NULL) |
| 3 | Conflict khi sync (đã abort, worktree giữ nguyên) |
| 4 | Bị chặn vì an toàn (dirty, commit chưa push) |
| 5 | Hết thời gian chờ khóa |

Ví dụ output `wtx ls`:

```
PROJECT kirindo   BASE origin/dev
NAME          BRANCH              ↑AHEAD ↓BEHIND  STATE   LAST COMMIT
worktree-18   feat/worktree-18    3      12      dirty   2h ago
worktree-19   feat/worktree-19    0      0       clean   5m ago
fix-otp       feat/fix-otp        5      40      gone    9d ago   ← nhánh remote đã xóa, có thể dọn
```

## 8. Luồng xử lý chính

Nguyên tắc chung: khóa chỉ bao quanh các thao tác ghi vào `.git` dùng chung (fetch, worktree add/remove, branch delete), chỉ vài giây. Agent làm việc trong worktree của mình không bao giờ bị khóa.

### 8.1 `wtx new worktree-18`

1. Validate tên; tính `branch = feat/worktree-18`, `path = <worktrees_root>/kirindo/worktree-18`.
2. Đọc `wtx.baseBranch`. Nếu NULL: có TTY thì hỏi (gợi ý từ remote) rồi lưu; không TTY thì thoát code 2.
3. Kiểm tra: path chưa tồn tại, nhánh chưa có local (`git show-ref`), chưa có trên remote (sau fetch).
4. **Lấy khóa** repo.
5. `git fetch origin dev --prune` — cập nhật `origin/dev` mới nhất.
6. `git worktree add --no-track -b feat/worktree-18 <path> origin/dev` — tạo thẳng từ remote, không đụng đến `dev` local. `--no-track` để upstream không trỏ vào `dev` (tránh lỡ `git push` lên dev).
7. **Nhả khóa.**
8. Chép `copy_files` từ repo chính sang (bỏ qua file không tồn tại, không ghi đè).
9. Chạy `postCreate` trong worktree mới (lỗi hook chỉ cảnh báo, không xóa worktree).
10. In đường dẫn (hoặc JSON). Lần push đầu: `git push -u origin feat/worktree-18`.

Rollback: nếu bước 6 lỗi giữa chừng thì `git worktree remove --force` + xóa nhánh vừa tạo.

### 8.2 `wtx sync worktree-18`

1. Worktree dirty: dừng với code 4, trừ khi có `--autostash`.
2. Khóa → `git fetch origin dev` → nhả khóa.
3. `git rev-list --left-right --count origin/dev...HEAD` → behind = 0 thì xong.
4. Dự đoán: `git merge-tree --write-tree origin/dev HEAD` — không đụng file, trả danh sách file conflict.
5. Có conflict: in danh sách, hỏi tiếp tục (TTY) hoặc thoát code 3 (`--dry-run` luôn dừng ở đây).
6. Thực hiện `git rebase origin/dev` hoặc `git merge origin/dev`.
7. Nếu conflict thật: mặc định `--abort` để trả về trạng thái cũ và báo code 3; `--keep-conflicts` để để lại cho người/agent giải quyết.
8. Rebase mà nhánh đã push: nhắc dùng `git push --force-with-lease` (tool không tự force push).

### 8.3 `wtx rm` (một hoặc nhiều worktree)

1. Xác định danh sách: theo tên, `--merged`, `--stale <ngày>`, hoặc chọn nhiều trong `-i`.
2. Với từng worktree, kiểm tra: không phải worktree chính; không dirty; không có commit chưa push; không bị `git worktree lock`; không có process đang đứng trong thư mục (Windows sẽ không cho xóa).
3. "Đã merge" được xác định theo thứ tự: `git merge-base --is-ancestor` (merge thường) → nhánh remote đã bị xóa sau `fetch --prune` (GitLab xóa nhánh sau squash merge) → v2: hỏi GitLab API trạng thái MR.
4. Hiển thị bảng tổng hợp (tên, lý do an toàn/chặn, dung lượng giải phóng) → xác nhận một lần.
5. Khóa → `git worktree remove <path>` → `git branch -d` (hoặc `-D` khi `--force`) → `git worktree prune` → nhả khóa.
6. Xóa thư mục dự án trong `worktrees/` nếu đã rỗng.

### 8.4 `wtx project list`

1. Liệt kê thư mục con cấp 1 của `projects_root`, bỏ `worktrees`.
2. Giữ thư mục có `.git` (thư mục hoặc file).
3. Với mỗi repo, đọc song song: `wtx.baseBranch` (NULL nếu trống), số worktree, remote URL.
4. In bảng. Không fetch — lệnh này phải nhanh và chạy được khi offline.

## 9. Rebase hay merge khi sync với dev

Khuyến nghị: **rebase khi nhánh chưa push hoặc chưa mở MR; merge khi MR đã mở và đang được review** — trừ khi team chấp nhận force push lên nhánh feature. Tool có thể tự chọn theo quy tắc này (`syncStrategy = auto`).

| Tiêu chí | Rebase | Merge |
| --- | --- | --- |
| Lịch sử | Thẳng, dễ review từng commit | Có merge commit "Merge origin/dev into…" |
| Nhánh đã push | Phải `push --force-with-lease`; reviewer mất vết comment trên commit cũ | Push bình thường |
| Nhiều commit + conflict | Giải conflict **mỗi commit một lần** (10 commit cùng sửa một file = có thể 10 lần) | Giải **một lần** |
| Agent tự giải conflict | Khó hơn, nhiều vòng | Dễ hơn |
| GitLab squash merge | Không quan trọng, cuối cùng cũng squash | Merge commit biến mất khi squash |

Về lo lắng "worktree có nhiều commit có thể gây xung đột":

- Luôn chạy dự đoán bằng `merge-tree` trước (bước 8.2.4). User/agent biết trước file nào sẽ conflict mà không làm bẩn worktree.
- Nhánh nhiều commit + có conflict: đề xuất merge, hoặc squash cục bộ trước rồi rebase (`git reset --soft $(git merge-base origin/dev HEAD)` + commit lại) — chỉ khi user đồng ý.
- Bật `git config rerere.enabled true` cho repo: git nhớ cách đã giải conflict và tự áp dụng lại khi gặp lại.
- Sync sớm và thường xuyên (`wtx ls` hiển thị số BEHIND lớn bằng màu cảnh báo) giảm conflict nhiều hơn mọi chiến lược.
- Không sync một worktree khi agent đang chạy trong đó: file bị thay dưới tay agent. MVP dựa vào kỷ luật; v2 có thể đánh dấu worktree "busy" bằng `git worktree lock --reason agent`.

## 10. Tích hợp agent, rủi ro, roadmap, câu hỏi mở

### Tích hợp với Claude agent

- Mỗi agent một worktree, khởi động `claude` từ trong thư mục worktree.
- Thêm vào `CLAUDE.md` của dự án: không `git checkout`/`git switch` sang nhánh khác, không sửa file ngoài worktree, dùng `wtx sync --dry-run` trước khi sync, không force push.
- `wtx` có thể in sẵn đoạn hướng dẫn đó: `wtx agent-rules >> CLAUDE.md` (v1).

### Rủi ro

| Rủi ro | Mức | Giảm thiểu |
| --- | --- | --- |
| Chạy git WSL và git Windows trên cùng repo | Cao | `doctor` cảnh báo; tài liệu hoá "mỗi repo một phía" |
| Xóa nhầm worktree còn việc | Cao | Kiểm tra dirty/unpushed, bảng xác nhận, `--force` mới bỏ qua |
| Windows không xóa được thư mục vì IDE/process giữ file | Trung bình | Báo rõ file bị khóa, cho chạy lại; `prune` dọn metadata |
| Dung lượng (`node_modules`, `bin/obj` × 4 worktree) | Trung bình | `wtx ls --size`, `rm --stale`; v2 cân nhắc pnpm store chung |
| Claude Code `-w` phát triển thêm và trùng tính năng | Thấp | Giữ trọng tâm vào base branch, multi-project, sync, cleanup |

### Roadmap

| Giai đoạn | Phạm vi | Ước lượng |
| --- | --- | --- |
| MVP | `init`, `project list/config`, `new`, `ls`, `rm <name>`, `path`, khóa repo, config 2 tầng | 1–2 tuần buổi tối |
| v1 | `sync` + dự đoán conflict, `rm --merged/--stale/-i`, `doctor`, `copy_files`, `postCreate`, `--json`, bản release 3 OS | +2 tuần |
| v2 | GitLab API (trạng thái MR, tạo MR), `--run claude`, đánh dấu busy, cấp port, `.wtx.toml` chia sẻ cho team, dashboard TUI | Tuỳ nhu cầu |

Ước lượng là cảm tính cho một dev làm ngoài giờ, có Claude hỗ trợ.

### Câu hỏi mở

- [ ] GitLab của công ty dùng squash merge hay merge commit? Có tự xóa nhánh nguồn sau khi merge không? (quyết định cách nhận biết "đã merge")
- [ ] Team có cho force push lên nhánh `feat/*` không? (quyết định rebase hay merge mặc định)
- [ ] 4 agent chạy trong WSL hay trên Windows? Repo hiện nằm ở `D:\` hay trong `~` của WSL?
- [ ] Có cần prefix khác ngoài `feat/` (ví dụ `fix/`, `hotfix/`) chọn theo từng task không?
- [ ] `claude --worktree` hiện tạo nhánh từ đâu và có cấu hình base được không — cần thử trước khi chốt phạm vi.

### Nguồn

- [Claude Code docs — Worktrees](https://code.claude.com/docs/en/worktrees)
- [Claude Code worktree support, v2.1.49 (codewithmukesh)](https://codewithmukesh.com/blog/git-worktrees-claude-code/)
