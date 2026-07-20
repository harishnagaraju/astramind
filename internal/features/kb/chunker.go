package kb

const (
	DefaultChunkSize = 1000
	DefaultOverlap   = 200
)

func ChunkDocument(doc *Document, chunkSize, overlap int) []Chunk {
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}

	if overlap < 0 {
		overlap = 0
	}

	var chunks []Chunk

	content := doc.Content
	index := 0
	start := 0

	for start < len(content) {

		end := start + chunkSize

		if end > len(content) {
			end = len(content)
		}

		chunks = append(chunks, Chunk{
			ID:         generateDocumentID(),
			DocumentID: doc.ID,
			Index:      index,
			Content:    content[start:end],
		})

		index++

		if end == len(content) {
			break
		}

		start = end - overlap
	}

	return chunks
}
