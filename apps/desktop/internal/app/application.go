package app

import (
	"context"

	"github.com/lsl1016/feedora/apps/desktop/internal/bridge"
)

type Application struct {
	workspace  *bridge.WorkspaceBridge
	repository *bridge.RepositoryBridge
	search     *bridge.SearchBridge
	remote     *bridge.RemoteBridge
	agent      *bridge.AgentBridge
	system     *bridge.SystemBridge
}

func New() *Application {
	return &Application{
		workspace:  bridge.NewWorkspaceBridge(),
		repository: bridge.NewRepositoryBridge(),
		search:     bridge.NewSearchBridge(),
		remote:     bridge.NewRemoteBridge(),
		agent:      bridge.NewAgentBridge(),
		system:     bridge.NewSystemBridge(),
	}
}

func (a *Application) Startup(ctx context.Context) { a.system.SetContext(ctx) }

func (a *Application) Bindings() []interface{} {
	return []interface{}{a.workspace, a.repository, a.search, a.remote, a.agent, a.system}
}
