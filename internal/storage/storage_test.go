package storage

import (
	"testing"
	"time"
)

func setupTestStorage() (*Storage, func()) {
	store, _ := NewStorage(5)
	cleanup := func() {
		store.Close()
	}
	return store, cleanup
}

func TestStorage_Add(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	clip := Clip{
		Content:   "test content",
		Type:      "text",
		Timestamp: time.Now(),
	}

	id, err := store.Add(clip)
	if err != nil {
		t.Fatalf("Failed to add clip: %v", err)
	}

	if id == 0 {
		t.Fatal("Expected non-zero ID")
	}
}

func TestStorage_Get(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	clip := Clip{
		Content:   "test content",
		Type:      "text",
		Timestamp: time.Now(),
	}

	id, _ := store.Add(clip)

	retrieved, err := store.Get(id)
	if err != nil {
		t.Fatalf("Failed to get clip: %v", err)
	}

	if retrieved.Content != clip.Content {
		t.Fatalf("Expected content %s, got %s", clip.Content, retrieved.Content)
	}
}

func TestStorage_List(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	// Add multiple clips
	for i := 0; i < 3; i++ {
		clip := Clip{
			Content:   "content " + string(rune('a'+i)),
			Type:      "text",
			Timestamp: time.Now(),
		}
		store.Add(clip)
	}

	clips, err := store.List(10)
	if err != nil {
		t.Fatalf("Failed to list clips: %v", err)
	}

	if len(clips) != 3 {
		t.Fatalf("Expected 3 clips, got %d", len(clips))
	}
}

func TestStorage_MemoryEviction(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	// Add 5 clips (exceeds memory limit of 5)
	for i := 0; i < 5; i++ {
		clip := Clip{
			Content:   "content " + string(rune('a'+i)),
			Type:      "text",
			Timestamp: time.Now(),
		}
		store.Add(clip)
	}

	// Should only keep last 5 clips in memory
	clips, err := store.List(10)
	if err != nil {
		t.Fatalf("Failed to list clips: %v", err)
	}

	if len(clips) != 5 {
		t.Fatalf("Expected 5 clips (memory limit), got %d", len(clips))
	}
}

func TestStorage_Remove(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	clip := Clip{
		Content:   "test content",
		Type:      "text",
		Timestamp: time.Now(),
	}

	id, _ := store.Add(clip)

	err := store.Remove(id)
	if err != nil {
		t.Fatalf("Failed to remove clip: %v", err)
	}

	_, err = store.Get(id)
	if err == nil {
		t.Fatal("Clip still exists after removal")
	}
}

func TestStorage_Clear(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	// Add clips
	for i := 0; i < 5; i++ {
		clip := Clip{
			Content:   "content " + string(rune('a'+i)),
			Type:      "text",
			Timestamp: time.Now(),
		}
		store.Add(clip)
	}

	// Clear all
	err := store.Clear()
	if err != nil {
		t.Fatalf("Failed to clear storage: %v", err)
	}

	// Verify everything is gone
	clips, err := store.List(10)
	if err != nil {
		t.Fatalf("Failed to list clips: %v", err)
	}

	if len(clips) != 0 {
		t.Fatalf("Expected 0 clips after clear, got %d", len(clips))
	}
}

func TestStorage_Search(t *testing.T) {
	store, cleanup := setupTestStorage()
	defer cleanup()

	// Add clips with different content
	clips := []Clip{
		{Content: "api_key_123", Type: "text", Timestamp: time.Now()},
		{Content: "database_url", Type: "text", Timestamp: time.Now()},
		{Content: "another value", Type: "text", Timestamp: time.Now()},
	}

	for _, clip := range clips {
		store.Add(clip)
	}

	// Search for "api"
	results, err := store.Search("api", 10)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("Expected search results")
	}

	// Should find "api_key_123"
	found := false
	for _, result := range results {
		if result.Content == "api_key_123" {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("Search didn't find 'api_key_123'")
	}
}
