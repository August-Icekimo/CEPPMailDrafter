# CEPPMailDrafter

CEPPMailDrafter 是一個以 Go 1.22+ 編寫的 CLI 工具，讀取本地目錄中預先存好的每月例行通知信件模板（Markdown + custom tags），將動態資料填入後，產生可直接供 Thunderbird 使用的 `.eml` 草稿檔案。

## 功能特性
- Markdown + YAML Front Matter 模板。
- 支援標籤：單純替換 `{{VAR}}`、條件區塊 `{{#IF_X}}...{{/IF_X}}`、迴圈區塊 `{{#LOOP}}...{{/LOOP}}`。
- YAML 資料補注（若檔案遺失將給予警告但能繼續提供預先定型的字串格式）。
- 支援 To, Cc, Bcc 陣列或單一字串。
- 支援本機附件 (`attachments:`)
- 中文主旨 Base64 Encoding
- `--dry-run` 模式，輸出至 Stdout 供預覽而不寫入檔案。
- 無外部相依 (僅使用了 `gopkg.in/yaml.v3`)。

## 環境需求
- Go 1.22+ 
- Thunderbird (若需檢視產生的草稿)

## 編譯與安裝
1. 取出或進入專案資料夾。
2. 執行 `go build -o maildraft ./cmd/maildraft` 完成建置。
3. 若需準備佈署至 Linux 主機 (HomeLab Linux box)，請透過交叉編譯即可：
   `GOOS=linux GOARCH=amd64 go build -o maildraft-linux ./cmd/maildraft`

## 使用教學

執行專案中內建的範例。請確保命令在專案根目錄下：
```bash
./maildraft --month example-2025-01 --dry-run
```
這會在螢幕上列印出產生的草稿預覽。

若有信心，可將其寫入為 `.eml` 草稿：
```bash
./maildraft --month example-2025-01
```
成功後，在 `output/` 目錄內即可找到 `example-2025-01.eml`，請將其拖曳至 Thunderbird。
