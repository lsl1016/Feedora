package bridge

import "errors"

var errWorkspaceNotWired = errors.New("workspace service is not wired yet")

type WorkspaceBridge struct{}

func NewWorkspaceBridge() *WorkspaceBridge { return &WorkspaceBridge{} }
func (b *WorkspaceBridge) Health() map[string]any { return map[string]any{"ready": false, "adapter": "local-workspace"} }
func (b *WorkspaceBridge) ListNotes(query map[string]any) (map[string]any, error) { return nil, errWorkspaceNotWired }
func (b *WorkspaceBridge) GetNote(noteID int) (map[string]any, error) { return nil, errWorkspaceNotWired }
func (b *WorkspaceBridge) CreateNote(input map[string]any) (map[string]any, error) { return nil, errWorkspaceNotWired }
func (b *WorkspaceBridge) UpdateNote(noteID int, input map[string]any) (map[string]any, error) { return nil, errWorkspaceNotWired }
func (b *WorkspaceBridge) DeleteNote(noteID int) error { return errWorkspaceNotWired }
func (b *WorkspaceBridge) ListKnowledgeBases() ([]map[string]any, error) { return nil, errWorkspaceNotWired }
