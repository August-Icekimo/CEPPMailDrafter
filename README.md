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

### CLI 參數說明

`maildraft` 提供以下參數供自訂路徑與執行模式：

- `--month <name>` (必填): 指定要處理的任務名稱或月份（對應 `${name}.md` 的模板與 `${name}.yaml` 的資料）。
- `--template-dir <path>`: 指定模板資料夾路徑，預設為 `./templates`。
- `--data-dir <path>`: 指定資料目錄路徑，預設為 `./data`。
- `--output-dir <path>`: 指定產生的 `.eml` 存放路徑，預設為 `./output`。
- `--dry-run`: 啟動預覽模式，不會產生實體檔案，而是將生成的 `.eml` 內容直接列印到終端機 (Stdout)。

### 單信件作法範例

執行專案中內建的範例。請確保命令在專案根目錄下：

1. **預覽草稿**
   ```bash
   ./maildraft --month example-2025-01 --dry-run
   ```
   這會在螢幕上列印出產生的草稿預覽，確認變數替換正確。

2. **產生實體信件**
   若預覽無誤，可將其寫入為 `.eml` 草稿：
   ```bash
   ./maildraft --month example-2025-01
   ```
   成功後，在 `output/` 目錄內即可找到 `example-2025-01.eml`，請將其拖曳至 Thunderbird 的草稿匣或任意資料夾中。

3. **搭配自訂目錄目錄**
   如果你將模板和資料放在自訂的目錄下，可以透過參數指定：
   ```bash
   ./maildraft --month my-newsletter --template-dir ./my-templates --data-dir ./my-data --output-dir ./send-ready
   ```

### 批次轉換作法範例

如果有多個月分或多封不同通報信件需要一次性產生，可以撰寫簡單的 Bash 迴圈來批次執行：

#### 方式 A：同時切換對應名稱的模板與資料
```bash
# 批次處理多個專案的當月通知信
for project in "project-A" "project-B" "project-C"; do
    ./maildraft --month "${project}-2025-01" 
done
```

#### 方式 B：將多個不同的模板，套用同一份資料檔
若需要將多個不同樣板（例如 `templates/CEPP*.md`）全部套用在同一個資料檔（例如 `data/CEPPWEB2026.yaml`），可以使用細分模式撰寫迴圈（可參考專案內的 `batch_cepp.sh`）：

```bash
# 指定共用的資料檔名稱
DATA_FILE="CEPPWEB2026.yaml"

# 迴圈處理 templates 目錄下所有 CEPP 開頭的 Markdown 檔案
for template_path in templates/CEPP*.md; do
    template_filename=$(basename "$template_path")
    
    # 使用細分模式 (方式 B) 分別指定模板與資料檔
    ./maildraft \
        --template-dir "./templates" \
        --template-file "$template_filename" \
        --data-dir "./data" \
        --data-file "$DATA_FILE"
done
```
成功後，`output/` 目錄內就會同時產生這些處理完成的 `.eml` 檔案。
