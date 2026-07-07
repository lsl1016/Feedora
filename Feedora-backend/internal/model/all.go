package model

// AllModels 返回全部模型，供 AutoMigrate 统一注册建表。
func AllModels() []interface{} {
	return []interface{}{
		// 阶段一核心表。
		&User{}, &Post{}, &PostImage{}, &Tag{}, &PostTag{}, &Topic{}, &PostTopic{},
		&Comment{}, &PostLike{}, &PostFavorite{}, &CommentLike{},
		&Circle{}, &CircleMember{}, &File{}, &OperationLog{},
		// 阶段二预留表（提前建表，为后续能力打基础）。
		&Notification{}, &UserPointLog{}, &EventOutbox{}, &WorkerEventRecord{},
	}
}
