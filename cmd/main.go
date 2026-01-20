package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/oscarhkli/statement-parser/cmd/internal/logging"
)

func main() {
	logging.Init()

	path := ""
	if len(os.Args) < 2 {
		log.Fatal("Please provide the PDF file path as an argument.")
	}
	path = os.Args[1]

	text, err := readPdfDirect(path)
	if err != nil {
		panic(err)
	}

	fmt.Println(text)

	fileName := strings.Split(path, ".")[0]
	err = writeFile(fileName+".txt", text)
	if err != nil {
		panic(err)
	}
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
		log.Fatalf("pdftotext failed: %v", err)
	}
	text := string(out)
	return text, nil
}
