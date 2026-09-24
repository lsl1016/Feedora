package bridge

type RemoteBridge struct{}

func NewRemoteBridge() *RemoteBridge { return &RemoteBridge{} }
func (b *RemoteBridge) Health() map[string]any { return map[string]any{"connected": false, "service": "feedora-server"} }
