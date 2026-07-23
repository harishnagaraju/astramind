package ai

import "github.com/harishnagaraju/astramind/internal/infrastructure/models"

type ChatRequest struct {
	Model    string
	APIKey   string
	Messages []models.Message

	// Temperature, if set, is passed through to the provider to
	// control generation randomness. A pointer so "unset" (nil) is
	// distinguishable from "explicitly zero" - callers that don't
	// care leave this nil and get the provider's own default,
	// rather than every caller silently getting temperature 0.
	//
	// Enumeration/extraction tasks (like /kb ask) benefit from a low
	// value: higher temperatures were observed, during v0.9.1
	// validation, to make the model's enumeration of matching
	// entries from a knowledge base less consistent run-to-run on
	// an otherwise-unchanged prompt.
	Temperature *float64
}
