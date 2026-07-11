package kb

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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
	return filepath.Join(
		s.documentsDir(),
		id+".json",
	)
}

func (s *JSONStorage) SaveDocument(doc *Document) error {
	if err := os.MkdirAll(s.documentsDir(), 0755); err != nil {
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

		name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		doc, err := s.LoadDocument(name)
		if err != nil {
			continue
		}

		documents = append(documents, *doc)
	}

	return documents, nil
}

func (s *JSONStorage) SaveChunks(chunks []Chunk) error {
	if len(chunks) == 0 {
		return nil
	}

	if err := os.MkdirAll(s.chunksDir(), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(chunks, "", "  ")
	if err != nil {
		return err
	}

	documentID := chunks[0].DocumentID

	return os.WriteFile(
		s.chunkPath(documentID),
		data,
		0644,
	)
}

func (s *JSONStorage) LoadChunks(documentID string) ([]Chunk, error) {
	data, err := os.ReadFile(s.chunkPath(documentID))
	if err != nil {
		if os.IsNotExist(err) {
			return []Chunk{}, nil
		}
		return nil, err
	}

	var chunks []Chunk

	if err := json.Unmarshal(data, &chunks); err != nil {
		return nil, err
	}

	return chunks, nil
}

func (s *JSONStorage) DeleteChunks(documentID string) error {
	err := os.Remove(s.chunkPath(documentID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}

func (s *JSONStorage) documentsDir() string {
	return filepath.Join(s.directory, "documents")
}

func (s *JSONStorage) chunksDir() string {
	return filepath.Join(s.directory, "chunks")
}

func (s *JSONStorage) chunkPath(documentID string) string {
	return filepath.Join(
		s.chunksDir(),
		documentID+".json",
	)
}
