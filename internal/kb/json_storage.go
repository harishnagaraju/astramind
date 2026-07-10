package kb

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// JSONStorage stores documents as JSON files.
type JSONStorage struct {
	directory string
}

// NewJSONStorage creates a new JSON storage backend.
func NewJSONStorage(directory string) *JSONStorage {
	return &JSONStorage{
		directory: directory,
	}
}

func (s *JSONStorage) documentPath(id string) string {
	return filepath.Join(s.directory, id+".json")
}

func (s *JSONStorage) SaveDocument(doc *Document) error {
	if err := os.MkdirAll(s.directory, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		s.documentPath(doc.ID),
		data,
		0644,
	)
}

func (s *JSONStorage) LoadDocument(id string) (*Document, error) {
	data, err := os.ReadFile(s.documentPath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}

	var doc Document

	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

func (s *JSONStorage) DeleteDocument(id string) error {
	err := os.Remove(s.documentPath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return ErrDocumentNotFound
		}
		return err
	}

	return nil
}

func (s *JSONStorage) ListDocuments() ([]Document, error) {
	var documents []Document

	files, err := os.ReadDir(s.directory)
	if err != nil {
		if os.IsNotExist(err) {
			return documents, nil
		}
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		doc, err := s.LoadDocument(name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
		if err != nil {
			continue
		}

		documents = append(documents, *doc)
	}

	return documents, nil
}
