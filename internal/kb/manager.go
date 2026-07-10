package kb

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

