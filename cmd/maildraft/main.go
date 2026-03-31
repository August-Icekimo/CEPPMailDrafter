package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/icekimo/CEPPMailDrafter/internal/data"
	"github.com/icekimo/CEPPMailDrafter/internal/loader"
	"github.com/icekimo/CEPPMailDrafter/internal/parser"
	"github.com/icekimo/CEPPMailDrafter/internal/renderer"
	"github.com/icekimo/CEPPMailDrafter/internal/writer"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	month := flag.String("month", "", "快捷模式：同時讀取 <template-dir>/<name>.md 與 <data-dir>/<name>.yaml")
	templateDir := flag.String("template-dir", "./templates", "模板目錄（--month 或細分模式皆適用）")
	templateFile := flag.String("template-file", "", "模板檔名（含副檔名，細分模式必填，例如 monthly.md）")
	dataDir := flag.String("data-dir", "./data", "資料目錄（--month 或細分模式皆適用）")
	dataFile := flag.String("data-file", "", "資料檔名（含副檔名，細分模式必填，例如 monthly.yaml）")
	outputDir := flag.String("output-dir", "./output", "輸出目錄，用來存放生成的 .eml 檔案")
	dryRun := flag.Bool("dry-run", false, "預覽模式：輸出至 stdout，不寫入檔案")

	flag.Usage = printUsage
	flag.Parse()

	// ── 決定資料來源 ────────────────────────────────────────────────────────────
	var resolvedTemplateDir, resolvedTemplateFile string
	var resolvedDataDir, resolvedDataFile string
	var outputName string // 用於命名輸出檔（不含副檔名）

	if *month != "" {
		// 模式 A：--month 快捷模式，同名 .md + .yaml
		resolvedTemplateDir = *templateDir
		resolvedTemplateFile = *month + ".md"
		resolvedDataDir = *dataDir
		resolvedDataFile = *month + ".yaml"
		outputName = *month
	} else if *templateFile != "" && *dataFile != "" {
		// 模式 B：細分模式，明確指定檔名
		resolvedTemplateDir = *templateDir
		resolvedTemplateFile = *templateFile
		resolvedDataDir = *dataDir
		resolvedDataFile = *dataFile
		// 輸出名稱取模板檔名去掉副檔名
		outputName = strings.TrimSuffix(*templateFile, filepath.Ext(*templateFile))
	} else {
		// 兩種模式都不滿足 → 報錯並顯示範例
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "錯誤：請提供資料來源參數，使用下列其中一種方式：")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  方式 A（快捷）：同時讀取同名 .md 模板與 .yaml 資料")
		fmt.Fprintln(os.Stderr, "    maildraft --month 2025-01 [--template-dir ./templates] [--data-dir ./data]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  方式 B（細分）：分別指定模板檔與資料檔")
		fmt.Fprintln(os.Stderr, "    maildraft --template-dir ./templates --template-file monthly.md \\")
		fmt.Fprintln(os.Stderr, "              --data-dir ./data --data-file monthly.yaml")
		fmt.Fprintln(os.Stderr, "")
		if *templateFile != "" && *dataFile == "" {
			return fmt.Errorf("使用細分模式時 --data-file 為必填（已提供 --template-file=%s）", *templateFile)
		}
		if *dataFile != "" && *templateFile == "" {
			return fmt.Errorf("使用細分模式時 --template-file 為必填（已提供 --data-file=%s）", *dataFile)
		}
		return fmt.Errorf("必須提供 --month 或同時提供 --template-file 與 --data-file")
	}

	// ── 載入模板 ────────────────────────────────────────────────────────────────
	fl := loader.NewFileLoader(resolvedTemplateDir)
	rawContent, err := fl.LoadFile(resolvedTemplateDir, resolvedTemplateFile)
	if err != nil {
		return fmt.Errorf("load template: %w", err)
	}

	frontYAML, bodyText, err := parser.SplitFrontMatter(rawContent)
	if err != nil {
		return fmt.Errorf("parse front matter: %w", err)
	}

	tokens, err := parser.Tokenise(bodyText)
	if err != nil {
		return fmt.Errorf("parse body tokens: %w", err)
	}

	// ── 載入資料 ────────────────────────────────────────────────────────────────
	ds, err := data.NewYAMLSourceFromFile(resolvedDataDir, resolvedDataFile)
	if err != nil {
		return fmt.Errorf("load data source: %w", err)
	}

	// ── 渲染與輸出 ──────────────────────────────────────────────────────────────
	r := renderer.New(ds)
	msg, err := r.Render(frontYAML, tokens)
	if err != nil {
		return fmt.Errorf("render message: %w", err)
	}

	var w writer.Writer
	if *dryRun {
		w = writer.NewDryRunWriter(os.Stdout)
	} else {
		w = writer.NewEMLWriter(*outputDir)
	}

	if err := w.Write(msg, outputName); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	if !*dryRun {
		fmt.Printf("成功生成郵件草稿 %s → %s\n", outputName, *outputDir)
	}

	return nil
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "maildraft — 從 Markdown 模板與 YAML 資料生成 Thunderbird .eml 草稿")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "用法（方式 A，快捷，同名模板+資料）：")
	fmt.Fprintln(os.Stderr, "  maildraft --month <name> [--template-dir <path>] [--data-dir <path>] [--output-dir <path>] [--dry-run]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "用法（方式 B，細分）：")
	fmt.Fprintln(os.Stderr, "  maildraft --template-dir <path> --template-file <filename>")
	fmt.Fprintln(os.Stderr, "            --data-dir <path> --data-file <filename>")
	fmt.Fprintln(os.Stderr, "            [--output-dir <path>] [--dry-run]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "旗標說明：")
	flag.PrintDefaults()
}
