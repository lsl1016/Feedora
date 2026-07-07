package event

// Kafka Topic 常量（不含前缀，前缀由配置 kafka.topicPrefix 指定）。
const (
	TopicUser        = "user.events"
	TopicPost        = "post.events"
	TopicComment     = "comment.events"
	TopicInteraction = "interaction.events"
	TopicCircle      = "circle.events"
)
