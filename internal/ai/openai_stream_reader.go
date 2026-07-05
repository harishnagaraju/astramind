package ai

import "io"

func (p *OpenAIProvider) readStream(
	body io.ReadCloser,
	stream *openAIStream,
) {
	defer body.Close()
	defer close(stream.events)
}