package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/oscarhkli/statement-parser/cmd/internal/logging"
	"github.com/oscarhkli/statement-parser/cmd/statementparse"
)

func main() {
	logging.Init()

	if err := run(); err != nil {
		slog.Error("application fail", "err", err)
		os.Exit(1)
	}
}

func run() error {
	outputType := ""
	flag.StringVar(&outputType, "o", "json", "Output format {json|csv}")
	flag.Parse()

	args := flag.Args()
	path := ""
	if len(args) == 0 {
		return errors.New("Please provide the path to the PDF statement")
	} else if len(args) > 1 {
		flag.Usage()
		return errors.New("too many arguments provided")
	}

	path = args[0]

	text, err := readPdfDirect(path)
	if err != nil {
		return err
	}

	statement := statementparse.Parse(text)
	outputText := ""
	if outputType == "json" {
		jsonStr, err := statement.ToJSON()
		if err != nil {
			return errors.New("Failed to convert statement to JSON: " + err.Error())
		}
		outputText = jsonStr
	} else if outputType == "csv" {
		csvStr, err := statement.ToCSV()
		if err != nil {
			return errors.New("Failed to convert statement to CSV: " + err.Error())
		}
		outputText = csvStr
	}

	fileName := strings.Split(path, ".")[0]
	return writeFile(fmt.Sprintf("%s.%s", fileName, outputType), outputText)
}

func writeFile(path string, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func readPdfDirect(path string) (string, error) {
	cmd := exec.Command("pdftotext", "-layout", "-nopgbrk", path, "-")
	out, err := cmd.Output()
	if err != nil {
		return "", errors.New("pdftotext failed: " + err.Error())
	}
	text := string(out)
	return text, nil
}
