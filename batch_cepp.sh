/.#!/bin/bash

# 指定共用的資料檔名稱
DATA_FILE="CEPPWEB2026.yaml"
DATA_PATH="./data/$DATA_FILE"

# 檢查資料檔中的年月是否與現在時間一致
if [[ ! -f "$DATA_PATH" ]]; then
    echo "找不到資料檔: $DATA_PATH"
    exit 1
fi

data_year=$(sed -n 's/^YEAR:[[:space:]]*"\{0,1\}\([0-9]\{4\}\)"\{0,1\}.*/\1/p' "$DATA_PATH")
data_month=$(sed -n 's/^MONTH:[[:space:]]*"\{0,1\}\([0-9]\{1,2\}\)"\{0,1\}.*/\1/p' "$DATA_PATH")
now_year=$(date +%Y)
now_month=$(date +%-m)

if [[ "$data_year" != "$now_year" || "$data_month" != "$now_month" ]]; then
    echo "注意：資料檔 $DATA_FILE 的年月為 ${data_year}/${data_month}，現在時間為 ${now_year}/${now_month}。"
    read -r -p "是否改用現在時間 ${now_year}/${now_month}？ [Y/n] " answer
    answer=${answer:-Y}
    if [[ "$answer" =~ ^[Yy]$ ]]; then
        sed -i "s/^YEAR:.*/YEAR: \"$now_year\"/" "$DATA_PATH"
        sed -i "s/^MONTH:.*/MONTH: \"$now_month\"/" "$DATA_PATH"
        echo "已將 $DATA_FILE 更新為 ${now_year}/${now_month}。"
    else
        echo "保留原設定 ${data_year}/${data_month}。"
    fi
fi

echo "開始批次產生 CEPP 信件..."

# 迴圈處理 templates 目錄下所有 CEPP 開頭的 Markdown 檔案
for template_path in templates/CEPP*.md; do
    # 確保檔案存在，避免找不到檔案時發生錯誤
    if [[ ! -f "$template_path" ]]; then
        echo "在 templates 目錄中找不到任何 CEPP*.md 檔案。"
        break
    fi
    
    # 取得不含路徑的檔名，例如 CEPP_A.md
    template_filename=$(basename "$template_path")
    
    echo "========================================"
    echo "正在處理模板: $template_filename"
    
    # 呼叫 maildraft，使用細分模式 (方式 B) 分別指定模板與資料檔
    ./maildraft \
        --template-dir "./templates" \
        --template-file "$template_filename" \
        --data-dir "./data" \
        --data-file "$DATA_FILE"
        
    if [[ $? -eq 0 ]]; then
        echo "成功處理: $template_filename"
    else
        echo "處理失敗: $template_filename"
    fi
done

echo "========================================"
echo "批次處理完成！請至 output 目錄檢查結果。"