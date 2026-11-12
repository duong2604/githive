---
id: REQ
title: Original requirements — CLI Tool For Worktree
status: Approved
req: []
depends_on: []
owner: Danny (duong2604)
last_updated: 13-11-2026
---

# Requirements — CLI Tool For Worktree

The text below is the original requirement written by Danny, kept verbatim in Vietnamese.
The only additions are the bold `[REQ-xx]` / `[NFR-xx]` tags and the sections after the horizontal rule.
Backslashes are escaped (`\\`) so the rendered text matches the original document.

## Original text

1. Ý tưởng tôi là một dev trong công ty chuyên lĩnh vực outsource. Tôi làm việc trên nhiều dự án sử dụng đồng thời nhiều claude agents cùng hỗ trợ làm việc. Hệ điều hành tôi đang sử dụng WSL. Nhận thấy một vấn đề trong công việc, khi tôi tham gia một dự án có tên X, được quản lý bởi gitlab server công ty, tôi thường sử dụng 4 claude agents chạy đồng thời để thực hiện các task khác nhau.
2. Khó khăn: 
   1. Nhánh làm việc chung "dev" luôn cập nhật liên tục bởi các thay đổi của dự án. Mỗi khi cần thực hiện một task mới tôi luôn cần pull code mới nhất từ nhánh này và tạo nhánh mới để làm việc. Khi xong task, tôi sẽ tạo pull request và chời review bởi team leader. 
   2. Kết thúc một task là khi được merge vào dev tôi cần pull lại code mới từ dev và tiếp tục tạo nhánh mới cho task. Tuy nhiên tại một thời điểm, chỉ một agent có thể làm việc với git, các agent còn lại phải chờ agent trước chạy xong mới có quyền checkout sang nhánh mới. 
3. **[NFR-01] [NFR-05]** Mong muốn tạo một công cụ cli tool tương thích với nhiều môi trường WSL, Ubuntu, Window, MAC OS và cho phép cài đặt và lựa chọn các thiết lập tiện ích thay vì gõ lệnh thủ công, đáp ứng các tiêu chí sau:
   1. **[NFR-02]** Có thể làm việc đồng thời trên nhiều nhánh cho từng task bởi claude agent 
   2. **[NFR-02]** Các agent không bị xung đột, mỗi agent chỉ làm việc trên một nhánh riêng của chính nó
   3. **[NFR-03]** Mỗi khi tạo một nhánh mới cho từng task luôn lấy thông tin commit mới nhất từ nhánh chỉ định. Trong trường hợp của tôi chúng tôi đặt tên nhánh là "dev", tuy nhiên mỗi dự án lại chỉ định tên khác nhau ví dụ "developer"
4. Tính năng giả định:
   1. **[REQ-01]** Lựa chọn đường dẫn thư mục làm việc khi cài đặt ban đầu, nếu chưa có thư mục có thể chọn thư mục đường dẫn mặc định với tên chỉ định mặc định từ tool. Ví dụ: trên window sẽ ở ổ D:\\\\ workspace\\\\worktrees\\\\ (mỗi task là một worktree con trong thư mục này). Lưu ý nếu trên window không có ổ D: có thể chọn ổ C: để tạo thư mục cha để làm việc. Tuy nhiên tôi chưa cân nhắc chính xác nên quyết định như nào? Quy tắc chọn đường dẫn thư mục làm việc như thế nào hoặc yêu cầu người dùng sử dụng đường dẫn mặc định do tôi quy định sẵn.
   2. **[REQ-02]** Cho phép lựa chọn thư mục chứa các dự án đang làm việc. Ví dụ: yêu cầu quy định chỉ ra đường dẫn đến thư mục cha chứa các dự án ở đây tôi giả sử D:\\\\workspace
   3. **[REQ-03]** Lựa chọn dự án đang làm việc để thiết lập làm việc tạo các worktree. Ví dụ D:\\\\workspace\\kirindo (mỗi dự án phải có .git)
   4. **[REQ-04]** Cho phép tạo worktree và đặt tên worktree. Khi tạo worktree ở bước này cần đảm bảo nhánh làm việc chính được khai báo, tôi gọi là staging hoặc dev, tôi giả định gọi là nhánh "dev" trên github / gitlab. Mỗi khi tạo mới một worktree cho phép user đặt tên worktree branch để làm việc. Khi người dùng đặt xong tên worktree giả định "worktree-18", tool sẽ tự động lấy commit mới nhất trên nhánh chính chỉ định - tôi giả định nhánh "dev" sau đó sẽ tạo worktree với tên chỉ định tương ứng trong đường dẫn chỉ định lưu trữ worktrees chung, đang giả định D:\\\\workspace\\\\worktrees\\\\worktree-18, tương ứng với tên nhánh bắt đầu bằng tiền tố "feat/" + tên worktree định nghĩa. Ví dụ ta có branch: feat/worktree-18
   5. **[REQ-05]** Cho phép tính năng hiển thị tất cả các worktree của một dự án 
   6. **[REQ-06]** Cho phép tính năng pull code mới nhất từ  nhánh “dev” và rebase / merge các thay đổi đó lên worktree đang làm việc. Tôi cần cân nhắc nếu worktree hiện tại đang có nhiều commit có thể gây xung đột cần tư vấn đánh giá thêm 
   7. **[REQ-07]** Tính năng xóa nhiều worktree của một dự án để làm sạch bộ nhớ máy tính và không gian lưu trữ khi làm việc xong với worktree cũ. Sau khi làm việc xong mà không xóa worktree vẫn tồn tại trong D:\\\\workspace\\\\worktrees\\\\
   8. **[REQ-08] [NFR-04]** Tính năng show danh sách dự án đang nằm trong workspace\\\\. Tính năng kiểm tra nhánh chỉ định chính để tạo workspace của từng dự án, nếu chưa  cấu hình đẻ NULL, nếu cấu hình khai báo thì hiển thị tên nhánh tương ứng. Đảm bảo các thông tin cài đặt không bị xóa hoặc mất khi hết phiên làm việc hoặc tắt máy.

---

## Requirement IDs

`BD-00-overview.md` refers to the feature items above as 4.1–4.8.

| ID | Summary | Source | BD |
| --- | --- | --- | --- |
| REQ-01 | Choose or initialise the worktrees folder, with an OS-specific default path | 4.1 | BD-01 |
| REQ-02 | Choose the parent folder that holds the projects | 4.2 | BD-01 |
| REQ-03 | Choose the working project (must contain `.git`) | 4.3 | BD-02 |
| REQ-04 | Create a worktree from the latest commit of the base branch, on branch `feat/<name>` | 4.4 | BD-04, BD-03 |
| REQ-05 | Show every worktree of a project | 4.5 | BD-05 |
| REQ-06 | Pull the base branch, then rebase/merge it into the worktree, warning about conflicts | 4.6 | BD-06 |
| REQ-07 | Remove several worktrees to free disk space | 4.7 | BD-07 |
| REQ-08 | List projects and each project's base branch (`NULL` when not configured) | 4.8 | BD-02 |
| NFR-01 | Runs on WSL, Ubuntu, Windows and macOS | 3 | BD-00, BD-09 |
| NFR-02 | Several agents work in parallel, one branch each, without conflicts | 3.1, 3.2 | BD-00, BD-04 |
| NFR-03 | The base branch name differs per project | 3.3 | BD-02 |
| NFR-04 | Configuration survives the end of a session and a reboot | 4.8 (last sentence) | BD-01 |
| NFR-05 | Install and choose settings through an interactive UI instead of typing commands | 3 | BD-01, BD-00 |

## Traceability

| REQ | BD | DD | Issue | Status |
| --- | --- | --- | --- | --- |
| REQ-01 | BD-01 | DD-01 | — | Not started |
| REQ-02 | BD-01 | DD-01 | — | Not started |
| REQ-03 | BD-02 | DD-02 | — | Not started |
| REQ-04 | BD-04, BD-03 | DD-04, DD-03 | — | Not started |
| REQ-05 | BD-05 | DD-05 | — | Not started |
| REQ-06 | BD-06 | DD-06 | — | Not started |
| REQ-07 | BD-07 | DD-07 | — | Not started |
| REQ-08 | BD-02 | DD-02 | — | Not started |
| NFR-01 | BD-00, BD-09 | DD-00, DD-09 | — | Not started |
| NFR-02 | BD-00, BD-04 | DD-00, DD-04 | — | Not started |
| NFR-03 | BD-02 | DD-02 | — | Not started |
| NFR-04 | BD-01 | DD-01 | — | Not started |
| NFR-05 | BD-01, BD-00 | DD-01, DD-00 | — | Not started |

**Not derived from a requirement:** BD-08 (`path`, `shell-init`) has no source requirement. It comes from the design, so that an agent or user can move into a worktree (`cd "$(githive path <target>)"`).

How to assign IDs and update this table: see [`docs/templates/REQ-format.md`](../templates/REQ-format.md).

## Source

Converted from `docs/Requirements/CLI Tool For Worktree.docx`, which was added in commit
`a28e746cb83e907ee56b4ded41d1d7d9121ad964` and removed after conversion. Retrieve the original with:

```
git show "a28e746:docs/Requirements/CLI Tool For Worktree.docx" > original.docx
```
