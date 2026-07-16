package models

import "testing"

func TestMessage(t *testing.T) {

	msg := Message{
		Role:    "user",
		Content: "Hello AstraMind",
	}

	if msg.Role != "user" {
		t.Fatal("Role not stored correctly")
	}

	if msg.Content != "Hello AstraMind" {
		t.Fatal("Content not stored correctly")
	}
}

func TestSearchResult(t *testing.T) {

	result := SearchResult{
		Index:   3,
		Role:    "assistant",
		Content: "Search result",
	}

	if result.Index != 3 {
		t.Fatal("Index not stored correctly")
	}

	if result.Role != "assistant" {
		t.Fatal("Role not stored correctly")
	}

	if result.Content != "Search result" {
		t.Fatal("Content not stored correctly")
	}
}

func TestSessionSearchResult(t *testing.T) {

	result := SessionSearchResult{
		Session: "default",
		Index:   7,
		Role:    "user",
		Content: "Session search",
	}

	if result.Session != "default" {
		t.Fatal("Session not stored correctly")
	}

	if result.Index != 7 {
		t.Fatal("Index not stored correctly")
	}

	if result.Role != "user" {
		t.Fatal("Role not stored correctly")
	}

	if result.Content != "Session search" {
		t.Fatal("Content not stored correctly")
	}
}
