package ai

type Provider interface {
	Name() string

	Chat(
		request ChatRequest,
	) (string, error)
}
