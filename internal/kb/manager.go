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

func (m *Manager) SaveChunks(chunks []Chunk) error {
	return m.storage.SaveChunks(chunks)
}

func (m *Manager) LoadDocument(id string) (*Document, error) {
	return m.storage.LoadDocument(id)
}

func (m *Manager) LoadChunks(documentID string) ([]Chunk, error) {
	return m.storage.LoadChunks(documentID)
}

func (m *Manager) DeleteDocument(id string) error {
	return m.storage.DeleteDocument(id)
}

func (m *Manager) DeleteChunks(documentID string) error {
	return m.storage.DeleteChunks(documentID)
}

// ListKnowledge returns all knowledge base documents.
func (m *Manager) ListKnowledge() ([]Document, error) {
	return m.ListDocuments()
}

// GetKnowledge returns a knowledge base document by ID.
func (m *Manager) GetKnowledge(documentID string) (*Document, error) {
	return m.LoadDocument(documentID)
}

// RemoveKnowledge removes a document and its chunks.
func (m *Manager) RemoveKnowledge(documentID string) error {

	if err := m.DeleteChunks(documentID); err != nil {
		return err
	}

	if err := m.DeleteDocument(documentID); err != nil {
		return err
	}

	return nil
}

// ClearKnowledge removes every knowledge base document.
func (m *Manager) ClearKnowledge() error {

	documents, err := m.ListKnowledge()
	if err != nil {
		return err
	}

	for _, doc := range documents {
		if err := m.RemoveKnowledge(doc.ID); err != nil {
			return err
		}
	}

	return nil
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
	}

	// Split the document into chunks.
	chunks := ChunkDocument(doc, DefaultChunkSize, DefaultOverlap)

	// Record the number of generated chunks.
	doc.ChunkCount = len(chunks)

	// Persist the chunks.
	if err := m.SaveChunks(chunks); err != nil {
		return nil, err
	}

	// Persist the document.
	if err := m.SaveDocument(doc); err != nil {
		return nil, err
	}

	return doc, nil
}
