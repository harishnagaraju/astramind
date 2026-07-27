package kb

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDeleteDocument_Success confirms the real, working case: a
// document that exists is actually removed from disk.
func TestDeleteDocument_Success(t *testing.T) {
	dir := t.TempDir()
	storage := NewJSONStorage(dir)

	doc := &Document{ID: "doc-1", Name: "test.txt", Content: "hello"}
	if err := storage.SaveDocument(doc); err != nil {
		t.Fatalf("setup: SaveDocument failed: %v", err)
	}

	if err := storage.DeleteDocument("doc-1"); err != nil {
		t.Fatalf("expected successful delete, got error: %v", err)
	}

	if _, err := storage.LoadDocument("doc-1"); err == nil {
		t.Fatal("expected document to be gone after delete, but it still loads")
	}
}

// TestDeleteDocument_NotFound confirms the documented, real behavior:
// deleting a document that never existed returns ErrDocumentNotFound,
// not a silent success and not a raw os.PathError.
func TestDeleteDocument_NotFound(t *testing.T) {
	dir := t.TempDir()
	storage := NewJSONStorage(dir)

	err := storage.DeleteDocument("never-existed")

	if err != ErrDocumentNotFound {
		t.Fatalf("expected ErrDocumentNotFound, got: %v", err)
	}
}

// TestDeleteDocument_GenuineOSError covers the real error path this
// coverage gap was actually about: a failure that is NOT "file
// doesn't exist" (e.g. a permissions problem, a locked file). Unix
// permission bits (chmod) aren't reliably testable across platforms
// (Windows doesn't honor them the same way), so this uses a
// portable, cross-platform-safe substitute: putting a non-empty
// directory where a file is expected. os.Remove genuinely fails on a
// non-empty directory on every OS, and critically, that failure is
// NOT os.IsNotExist - it exercises the exact "return err" branch
// (the real error, not the not-found branch) that was previously
// uncovered.
func TestDeleteDocument_GenuineOSError(t *testing.T) {
	dir := t.TempDir()
	storage := NewJSONStorage(dir)

	docPath := filepath.Join(dir, "documents", "blocked.json")
	if err := os.MkdirAll(docPath, 0755); err != nil {
		t.Fatalf("setup: failed to create blocking directory: %v", err)
	}
	// A non-empty directory in place of the expected file - os.Remove
	// fails on this, portably, on every OS.
	if err := os.WriteFile(filepath.Join(docPath, "occupied.txt"), []byte("x"), 0644); err != nil {
		t.Fatalf("setup: failed to occupy blocking directory: %v", err)
	}

	err := storage.DeleteDocument("blocked")

	if err == nil {
		t.Fatal("expected a genuine OS error (non-empty directory in place of file), got nil")
	}
	if err == ErrDocumentNotFound {
		t.Fatal("expected a genuine OS error, not ErrDocumentNotFound - the path does exist, just as an unexpected directory")
	}
}

// TestDeleteChunks_Success confirms the real, working case.
func TestDeleteChunks_Success(t *testing.T) {
	dir := t.TempDir()
	storage := NewJSONStorage(dir)

	chunks := []Chunk{{ID: "c1", DocumentID: "doc-1", Content: "hello"}}
	if err := storage.SaveChunks(chunks); err != nil {
		t.Fatalf("setup: SaveChunks failed: %v", err)
	}

	if err := storage.DeleteChunks("doc-1"); err != nil {
		t.Fatalf("expected successful delete, got error: %v", err)
	}

	loaded, err := storage.LoadChunks("doc-1")
	if err != nil {
		t.Fatalf("LoadChunks after delete should not error, got: %v", err)
	}
	if len(loaded) != 0 {
		t.Fatalf("expected no chunks after delete, got %d", len(loaded))
	}
}

// TestDeleteChunks_NotFoundReturnsNilNotError documents the real,
// confirmed asymmetry with DeleteDocument: deleting chunks for a
// document that was never chunked (or already had its chunks
// deleted) returns nil - success - not an error. DeleteDocument
// treats the identical "file doesn't exist" condition as an error
// (ErrDocumentNotFound); DeleteChunks treats it as already-satisfied.
// This may be intentional (a document with zero chunks is a valid,
// unremarkable state) but it was previously an unverified assumption,
// not a tested contract - this test makes the actual behavior
// explicit and will fail loudly if it's ever changed unintentionally.
func TestDeleteChunks_NotFoundReturnsNilNotError(t *testing.T) {
	dir := t.TempDir()
	storage := NewJSONStorage(dir)

	err := storage.DeleteChunks("never-had-chunks")

	if err != nil {
		t.Fatalf("expected nil (documented behavior: missing chunks is not an error), got: %v", err)
	}
}

// TestDeleteChunks_GenuineOSError mirrors the DeleteDocument OS-error
// test, for the same real, previously-uncovered branch in
// DeleteChunks.
func TestDeleteChunks_GenuineOSError(t *testing.T) {
	dir := t.TempDir()
	storage := NewJSONStorage(dir)

	chunkPath := filepath.Join(dir, "chunks", "blocked.json")
	if err := os.MkdirAll(chunkPath, 0755); err != nil {
		t.Fatalf("setup: failed to create blocking directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(chunkPath, "occupied.txt"), []byte("x"), 0644); err != nil {
		t.Fatalf("setup: failed to occupy blocking directory: %v", err)
	}

	err := storage.DeleteChunks("blocked")

	if err == nil {
		t.Fatal("expected a genuine OS error (non-empty directory in place of file), got nil")
	}
}
