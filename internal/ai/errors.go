package ai

import (
	"fmt"
	"net/http"
)

func handleAPIError(
	status int,
	body string,
) error {

	switch status {

	case http.StatusUnauthorized:
		return fmt.Errorf(
			"invalid API key",
		)

	case http.StatusPaymentRequired:
		return fmt.Errorf(
			"provider account requires payment",
		)

	case http.StatusForbidden:
		return fmt.Errorf(
			"access denied",
		)

	case http.StatusTooManyRequests:
		return fmt.Errorf(
			"rate limit or quota exceeded",
		)

	default:
		return fmt.Errorf(
			"API Error (%d): %s",
			status,
			body,
		)
	}
}
