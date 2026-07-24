[CmdletBinding()]
param (
    [string]$ExeName = "",
    [string]$DataFile = "CEPPWEB2026.yaml"
)

$DataPath = Join-Path "./data" $DataFile

# 檢查資料檔是否存在
if (-not (Test-Path $DataPath)) {
    Write-Host "找不到資料檔: $DataPath" -ForegroundColor Red
    exit 1
}

# 讀取檔案內容
$content = Get-Content -Path $DataPath -Raw -Encoding UTF8

# 檢查資料檔中的年月是否與現在時間一致
$dataYear = if ($content -match '(?m)^YEAR:\s*"?(\d{4})"?') { $Matches[1] } else { "" }
$dataMonth = if ($content -match '(?m)^MONTH:\s*"?(\d{1,2})"?') { $Matches[1] } else { "" }

$nowYear = (Get-Date).Year.ToString()
$nowMonth = (Get-Date).Month.ToString()

if ($dataYear -ne $nowYear -or $dataMonth -ne $nowMonth) {
    Write-Host "注意：資料檔 $DataFile 的年月為 ${dataYear}/${dataMonth}，現在時間為 ${nowYear}/${nowMonth}。" -ForegroundColor Yellow
    $answer = Read-Host "是否改用現在時間 ${nowYear}/${nowMonth}？ [Y/n]"
    if ([string]::IsNullOrWhiteSpace($answer)) { $answer = "Y" }

    if ($answer -match '^[Yy]$') {
        $content = $content -replace '(?m)^YEAR:.*', "YEAR: `"$nowYear`""
        $content = $content -replace '(?m)^MONTH:.*', "MONTH: `"$nowMonth`""
        Set-Content -Path $DataPath -Value $content -Encoding UTF8
        Write-Host "已將 $DataFile 更新為 ${nowYear}/${nowMonth}。" -ForegroundColor Green
    } else {
        Write-Host "保留原設定 ${dataYear}/${dataMonth}。"
    }
}

# 判斷執行檔名稱 (傳入參數 > maildraft-windows-amd64.exe > maildraft.exe > maildraft)
$exe = ""
if (-not [string]::IsNullOrWhiteSpace($ExeName)) {
    if (Test-Path ".\$ExeName") {
        $exe = ".\$ExeName"
    } elseif (Test-Path $ExeName) {
        $exe = $ExeName
    } else {
        Write-Host "找不到指定的執行檔: $ExeName" -ForegroundColor Red
        exit 1
    }
} else {
    $candidates = @(".\maildraft-windows-amd64.exe", ".\maildraft.exe", ".\maildraft")
    foreach ($c in $candidates) {
        if (Test-Path $c) {
            $exe = $c
            break
        }
    }
    if ([string]::IsNullOrWhiteSpace($exe)) {
        $exe = ".\maildraft-windows-amd64.exe"
    }
}

if (-not (Test-Path $exe)) {
    Write-Host "找不到執行檔 $exe ，請確認執行檔已放置於目前目錄下，或透過 -ExeName 參數指定檔名。" -ForegroundColor Red
    exit 1
}

Write-Host "使用執行檔: $exe" -ForegroundColor Cyan
Write-Host "開始批次產生 CEPP 信件..."

# 取得 templates 目錄下所有 CEPP 開頭的 Markdown 檔案
$templates = Get-ChildItem -Path "./templates" -Filter "CEPP*.md" -File

if ($null -eq $templates -or $templates.Count -eq 0) {
    Write-Host "在 templates 目錄中找不到任何 CEPP*.md 檔案。" -ForegroundColor Red
    exit 0
}

foreach ($tmpl in $templates) {
    $templateFilename = $tmpl.Name

    Write-Host "========================================"
    Write-Host "正在處理模板: $templateFilename"

    & $exe --template-dir "./templates" --template-file "$templateFilename" --data-dir "./data" --data-file "$DataFile"

    if ($LASTEXITCODE -eq 0) {
        Write-Host "成功處理: $templateFilename" -ForegroundColor Green
    } else {
        Write-Host "處理失敗: $templateFilename" -ForegroundColor Red
    }
}

Write-Host "========================================"
Write-Host "批次處理完成！請至 output 目錄檢查結果。" -ForegroundColor Cyan
