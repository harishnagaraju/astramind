package app

type App struct {
	deps         Dependencies
	dispatcher   *commandDispatcher
	apiKey       string
	model        string
	baseURL      string
	providerName string

	activeSession string
}

func New() *App {
	return &App{
		activeSession: "default",
	}
}
