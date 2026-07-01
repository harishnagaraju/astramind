package ai

import "sync"

var (
	mockCache []MockPair
	loadOnce  sync.Once
	loadError error
)

func InitializeMockCache() error {

	loadOnce.Do(func() {

		mockCache, loadError = LoadMockData()

	})

	return loadError
}
