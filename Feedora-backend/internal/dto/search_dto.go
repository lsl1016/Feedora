package dto

// SuggestItem 搜索联想项。
type SuggestItem struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	TargetID  int64  `json:"targetId"`
	TargetURL string `json:"targetUrl"`
}
