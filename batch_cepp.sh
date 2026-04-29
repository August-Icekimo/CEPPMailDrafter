/.#!/bin/bash

# 指定共用的資料檔名稱
DATA_FILE="CEPPWEB2026.yaml"

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