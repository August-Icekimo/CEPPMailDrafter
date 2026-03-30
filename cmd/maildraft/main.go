package main

import (
	"flag"
	"fmt"
	"os"

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
	month := flag.String("month", "", "Target month (e.g. 2025-01) - required")
	templateDir := flag.String("template-dir", "./templates", "Directory containing template files")
	dataDir := flag.String("data-dir", "./data", "Directory containing yaml data files")
	outputDir := flag.String("output-dir", "./output", "Directory to write the resulting eml file")
	dryRun := flag.Bool("dry-run", false, "Preview to stdout instead of writing an eml file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of maildraft:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *month == "" {
		flag.Usage()
		return fmt.Errorf("--month flag is required")
	}

	fl := loader.NewFileLoader(*templateDir)
	rawContent, err := fl.Load(*month)
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

	ds, err := data.NewYAMLSource(*month, *dataDir)
	if err != nil {
		return fmt.Errorf("load data source: %w", err)
	}

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

	if err := w.Write(msg, *month); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	if !*dryRun {
		fmt.Printf("Successfully generated mail draft for %s in %s\n", *month, *outputDir)
	}

	return nil
}
