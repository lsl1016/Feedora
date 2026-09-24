package bridge

import "errors"

var errRepositoryNotWired = errors.New("repository service is not wired yet")

type RepositoryBridge struct{}

func NewRepositoryBridge() *RepositoryBridge { return &RepositoryBridge{} }
func (b *RepositoryBridge) Health() map[string]any { return map[string]any{"ready": false, "adapter": "git-local"} }
func (b *RepositoryBridge) ListRepositories() ([]map[string]any, error) { return nil, errRepositoryNotWired }
func (b *RepositoryBridge) AddRepository(path string) (map[string]any, error) { return nil, errRepositoryNotWired }
func (b *RepositoryBridge) RemoveRepository(repositoryID string) error { return errRepositoryNotWired }
func (b *RepositoryBridge) ReadFile(repositoryID, path string) (map[string]any, error) { return nil, errRepositoryNotWired }
func (b *RepositoryBridge) Reindex(repositoryID string) error { return errRepositoryNotWired }
