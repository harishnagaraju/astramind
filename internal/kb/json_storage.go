package kb

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
