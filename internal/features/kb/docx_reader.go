package kb

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// extractDocxText extracts the visible paragraph text from a .docx
// file's raw bytes, with no external dependency: a .docx is a ZIP
// archive containing XML, and the standard library's archive/zip and
// encoding/xml packages are sufficient to read it directly.
//
// Paragraphs are joined with a blank line ("\n\n") between them,
// matching the convention chunker.go's splitParagraphs already
// expects - this keeps imported .docx content flowing through the
// exact same chunking, retrieval, and deterministic-extraction
// pipeline validated against real .txt content earlier, with no
// special-casing needed downstream for this format.
//
// Only word/document.xml (the main document body) is read. Headers,
// footers, footnotes, and embedded objects are not extracted - a
// deliberate first-pass scope limit, not an oversight; most contract/
// policy-style documents keep their substantive content in the body.
func extractDocxText(data []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("not a valid .docx file (not a zip archive): %w", err)
	}

	var documentXML *zip.File
	for _, f := range reader.File {
		if f.Name == "word/document.xml" {
			documentXML = f
			break
		}
	}

	if documentXML == nil {
		return "", fmt.Errorf("not a valid .docx file: word/document.xml not found")
	}

	rc, err := documentXML.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open word/document.xml: %w", err)
	}
	defer rc.Close()

	decoder := xml.NewDecoder(rc)

	var result strings.Builder
	var paragraph strings.Builder
	inText := false

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to parse word/document.xml: %w", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "t":
				// <w:t> holds a run of visible text. Namespace
				// prefixes ("w:") are already resolved by
				// encoding/xml, so checking Name.Local alone
				// (without the prefix) is correct here.
				inText = true
			case "tab":
				paragraph.WriteString("\t")
			case "br", "cr":
				// Explicit line break within a paragraph - not a
				// new paragraph, so a plain newline, not a blank
				// line.
				paragraph.WriteString("\n")
			}
		case xml.CharData:
			if inText {
				paragraph.Write(t)
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "t":
				inText = false
			case "p":
				text := strings.TrimSpace(paragraph.String())
				if text != "" {
					result.WriteString(text)
					result.WriteString("\n\n")
				}
				paragraph.Reset()
			}
		}
	}

	return strings.TrimSpace(result.String()), nil
}
