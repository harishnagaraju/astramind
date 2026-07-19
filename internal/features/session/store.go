package session

// Store is what Service needs for session lifecycle operations
// (Create/Exists) that don't belong in History. Satisfied
// structurally by *storage.FileHistoryStore - defined here so
// session doesn't depend on the concrete storage type, and so tests
// can substitute a store rooted at a temp directory instead of the
// real data/sessions folder.
type Store interface {
	CreateSession(session string) error
	SessionExists(session string) bool
}
