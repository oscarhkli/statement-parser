# statement-parser

Automate the boring stuff.

`statement-parser` extracts transactions from a PDF bank statement and outputs them as JSON or CSV files.

Currently, it only supports:

- HSBC (HK) Visa Signature statements  
- HSBC (HK) Red statements

## Prerequisites

- [`pdftotext`](https://poppler.freedesktop.org/) must be installed on your system.  
  - On macOS: `brew install poppler`  
  - On Ubuntu/Debian: `sudo apt install poppler-utils`  

`pdftotext` is used internally to convert PDF statements into text before parsing.

---

## Usage

Build the project and run the CLI:

```bash
make build
./bin/statement-parser -output={csv|json} <PDF_FILE>

Example:

```bash
./bin/statement-parser -output=json ~/Downloads/2025-10-20_Statement.pdf
```

[Visit oscarhkli.com for more](https://oscarhkli.com/)
