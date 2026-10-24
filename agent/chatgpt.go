package agent

import (
	"context"

	"github.com/Autumn-27/norma/llm"
)

// ChatGPT subscription requests must stream even when an internal caller needs
// a complete message (for example compaction). Reuse the SDK accumulator.
type chatGPTProvider struct{ llm.Provider }

func (p chatGPTProvider) Complete(ctx context.Context, req llm.CompletionRequest) (llm.Message, string, llm.Usage, error) {
	a := llm.NewAccumulator()
	for ev, err := range p.Provider.Stream(ctx, req) {
		if err != nil {
			return llm.Message{}, "", llm.Usage{}, err
		}
		a.Add(ev)
	}
	return a.Message(), a.StopReason, a.Usage, nil
}
