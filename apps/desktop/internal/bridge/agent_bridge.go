package bridge

import "errors"

var errAgentNotWired = errors.New("agent runtime adapter is not wired yet")

type AgentBridge struct{}

func NewAgentBridge() *AgentBridge { return &AgentBridge{} }
func (b *AgentBridge) Health() map[string]any { return map[string]any{"connected": false, "protocol": "runtime-adapter"} }
func (b *AgentBridge) CreateSession(title string) (string, error) { return "", errAgentNotWired }
func (b *AgentBridge) SendMessage(sessionID, message string, context []map[string]any) ([]map[string]any, error) {
	return nil, errAgentNotWired
}
func (b *AgentBridge) Cancel(sessionID string) error { return errAgentNotWired }
