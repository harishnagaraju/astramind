package kb

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manager coordinates knowledge base operations.
type Manager struct {
	storage Storage
}

// NewManager creates a new knowledge base manager.
func NewManager(storage Storage) *Manager {
	return &Manager{
		storage: storage,
	}
}

func (m *Manager) SaveDocument(doc *Document) error {
	return m.storage.SaveDocument(doc)
}

func (m *Manager) LoadDocument(id string) (*Document, error) {
	return m.storage.LoadDocument(id)
}

func (m *Manager) DeleteDocument(id string) error {
	return m.storage.DeleteDocument(id)
}

func (m *Manager) ListDocuments() ([]Document, error) {
	return m.storage.ListDocuments()
}

// ImportDocument imports a text or markdown file into the knowledge base.
func (m *Manager) ImportDocument(path string) (*Document, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".txt", ".md":
		// supported
	default:
		return nil, ErrInvalidDocument
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	id := generateDocumentID()

	doc := &Document{
		ID:         id,
		Name:       filepath.Base(path),
		Path:       path,
		Content:    string(data),
		CreatedAt:  time.Now(),
		ModifiedAt: info.ModTime(),
		ChunkCount: 0,
	}

	if err := m.SaveDocument(doc); err != nil {
		return nil, err
	}

	return doc, nil
}
