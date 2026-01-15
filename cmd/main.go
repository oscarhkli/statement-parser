package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	path := "2025-10-20_Statement"
	text, err := readPdfDirect(path + ".pdf")
	if err != nil {
		panic(err)
	}

	fmt.Println(text)
	err = writeFile(path+".txt", text)
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
