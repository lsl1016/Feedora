package bridge

import "errors"

var errSearchNotWired = errors.New("search service is not wired yet")

type SearchBridge struct{}

func NewSearchBridge() *SearchBridge { return &SearchBridge{} }
func (b *SearchBridge) Health() map[string]any { return map[string]any{"ready": false, "engine": "fts5"} }
func (b *SearchBridge) Search(query map[string]any) ([]map[string]any, error) { return nil, errSearchNotWired }
