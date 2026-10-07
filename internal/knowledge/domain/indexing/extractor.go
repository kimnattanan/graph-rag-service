package indexing

type ExtractResult struct {
	content  string
	entities []string
}

func NewExtractResult(content string, entities []string) ExtractResult {
	return ExtractResult{
		content:  content,
		entities: append([]string(nil), entities...),
	}
}

func (r *ExtractResult) Content() string {
	return r.content
}

func (r *ExtractResult) Entities() []string {
	return append([]string(nil), r.entities...)
}
