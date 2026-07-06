export type ID = number;
export type TargetType = 'post' | 'comment' | 'user' | 'tag' | 'circle' | 'topic' | 'activity' | 'announcement' | 'note' | 'knowledgeBase';
export type UserStatus = 'normal' | 'muted' | 'banned' | 'deleted';
export type UserRole = 'user' | 'admin';
export type SortType = 'recommend' | 'latest' | 'hot' | 'comment' | 'favorite' | 'view';
export type TimeRange = 'all' | 'today' | 'week' | 'month';
export type EnabledStatus = 'enabled' | 'disabled';

export interface PageRequest { page: number; pageSize: number; }
export interface PageResult<T> { list: T[]; total: number; page: number; pageSize: number; }
export interface UserSummary { userId: number; nickname: string; avatar: string; bio?: string; level?: number; levelName?: string; }
export interface User extends UserSummary {
  account: string; role: UserRole; status: UserStatus; points: number; experience: number; nextLevelExperience: number; badgeCount: number; checkedInToday: boolean; continuousCheckInDays: number;
  postCount: number; commentCount: number; followerCount: number; followingCount: number; likeReceivedCount: number; createdAt: string; updatedAt: string;
}
export interface LoginResult { token: string; user: User; }

export interface ContentTag { tagId: number; tagName: string; description: string; status: EnabledStatus; useCount: number; createdAt: string; updatedAt: string; }
export interface TagSummary { tagId: number; tagName: string; }
export interface Topic { topicId: number; name: string; description: string; coverImage: string; postCount: number; participantCount: number; isOfficial: boolean; isRecommended: boolean; status: EnabledStatus; createdAt: string; updatedAt: string; }
export interface TopicSummary { topicId: number; name: string; }

export type CircleStatus = 'normal' | 'reviewing' | 'closed' | 'deleted';
export type CircleJoinType = 'direct' | 'approval';
export type CirclePostPermission = 'all' | 'admin_only';
export type CircleMemberRole = 'owner' | 'moderator' | 'reviewer' | 'member';
export type CircleMemberStatus = 'normal' | 'muted' | 'removed';
export type CircleScopeTab = 'all' | 'recommended' | 'joined' | 'created';
export interface Circle { circleId: number; name: string; avatar: string; description: string; category: string; tags: TagSummary[]; ownerId: number; owner: UserSummary; joinType: CircleJoinType; postPermission: CirclePostPermission; memberCount: number; postCount: number; featuredPostCount: number; isJoined: boolean; myRole?: CircleMemberRole; myStatus?: CircleMemberStatus; isRecommended: boolean; status: CircleStatus; rules: string; createdAt: string; updatedAt: string; }
export interface CircleBrief { circleId: number; name: string; avatar: string; memberCount: number; postCount: number; }
export interface CircleMember { id: number; circleId: number; userId: number; user: UserSummary; role: CircleMemberRole; status: CircleMemberStatus; muteReason?: string; mutedUntil?: string; joinedAt: string; updatedAt: string; }
export interface CircleJoinRequest { requestId: number; circleId: number; userId: number; user: UserSummary; reason: string; status: 'pending' | 'approved' | 'rejected'; createdAt: string; handledAt?: string; handleReason?: string; }
export interface CircleAnnouncement { announcementId: number; circleId: number; title: string; content: string; publisher: UserSummary; createdAt: string; }

export type PostStatus = 'draft' | 'scheduled' | 'reviewing' | 'published' | 'hidden' | 'rejected' | 'deleted' | 'takedown';
export type PostType = 'original' | 'repost';
export type PostVisibility = 'public' | 'circle_only';
export interface Post { postId: number; postType: PostType; authorId: number; author: UserSummary; title: string; content: string; summary: string; images: string[]; tags: TagSummary[]; topics: TopicSummary[]; circleId?: number; circle?: CircleBrief; visibility: PostVisibility; status: PostStatus; isTop: boolean; isFeatured: boolean; isSelected: boolean; scheduledAt?: string; sourcePostId?: number; sourcePost?: Post; repostComment?: string; viewCount: number; likeCount: number; commentCount: number; favoriteCount: number; shareCount: number; repostCount: number; hotScore: number; liked: boolean; favorited: boolean; followedAuthor: boolean; createdAt: string; publishedAt?: string; updatedAt: string; }
export interface PostQuery extends PageRequest { feedType?: 'recommend' | 'latest' | 'hot' | 'following'; sort?: SortType; timeRange?: TimeRange; status?: PostStatus | 'all'; keyword?: string; tagId?: number; circleId?: number; topicId?: number; authorId?: number; includeHidden?: boolean; }
export interface CreatePostRequest { title: string; content: string; images: string[]; tagIds: number[]; topicIds: number[]; circleId?: number; visibility: PostVisibility; publishMode: 'now' | 'schedule' | 'draft'; scheduledAt?: string; }

export interface Comment { commentId: number; postId: number; userId: number; user: UserSummary; content: string; likeCount: number; liked: boolean; status: 'normal' | 'deleted' | 'rejected'; replies: Comment[]; createdAt: string; updatedAt: string; }
export interface MyCommentItem { commentId: number; postId: number; postTitle: string; content: string; likeCount: number; status: 'normal' | 'deleted' | 'rejected'; createdAt: string; }

export interface CircleChatMessage { messageId: number; circleId: number; senderId: number; sender: UserSummary; content: string; messageType: 'text' | 'system'; status: 'normal' | 'deleted'; createdAt: string; }
export type CircleSummaryFrequency = 'once' | 'daily' | 'weekly' | 'monthly';
export type CircleSummaryContentScope = 'post' | 'comment' | 'chat' | 'featured';
export interface CircleSummarySubscription { subscriptionId: number; userId: number; circleId: number; circleName?: string; enabled: boolean; frequency: CircleSummaryFrequency; pushTime?: string; oneTimePushAt?: string; channels: Array<'notification' | 'email'>; email?: string; contentScopes: CircleSummaryContentScope[]; createdAt: string; updatedAt: string; }

export type FollowingObjectTab = 'users' | 'circles' | 'topics' | 'tags';
export type FollowingFeedTab = 'all' | 'user' | 'circle' | 'topic';
export interface FollowingFeedItem { feedId: number; feedType: 'post' | 'repost' | 'circle' | 'topic' | 'tag'; title: string; summary: string; targetId: number; targetUrl: string; sourceName: string; sourceAvatar?: string; createdAt: string; }

export interface WorkspaceNote { noteId: number; title: string; content: string; summary: string; knowledgeBaseId?: number; knowledgeBaseName?: string; ownerId: number; status: 'normal' | 'deleted'; createdAt: string; updatedAt: string; }
export interface KnowledgeBase { knowledgeBaseId: number; name: string; description: string; ownerId: number; noteCount: number; memberCount: number; visibility: 'private' | 'collaborative'; createdAt: string; updatedAt: string; }
export interface KnowledgeBaseMember { id: number; knowledgeBaseId: number; userId: number; user: UserSummary; role: 'owner' | 'editor' | 'viewer'; status: 'active' | 'pending'; invitedAt?: string; joinedAt?: string; }
export interface WorkspaceDashboardStats { noteCount: number; knowledgeBaseCount: number; aiChatCount: number; collaborationMemberCount: number; }

export type AiProvider = 'openai' | 'deepseek' | 'qwen' | 'siliconflow' | 'custom';
export interface UserModelConfig { configId: number; userId: number; provider: AiProvider; baseUrl: string; apiKeyMasked: string; modelName: string; temperature: number; maxTokens: number; status: EnabledStatus; createdAt: string; updatedAt: string; }
export interface AiChatSession { sessionId: number; title: string; userId: number; createdAt: string; updatedAt: string; }
export interface AiChatMessage { messageId: number; sessionId: number; role: 'user' | 'assistant'; content: string; editable: boolean; createdAt: string; }

export interface Badge { badgeId: number; name: string; icon: string; description: string; condition: string; status: EnabledStatus; createdAt: string; }
export interface Task { taskId: number; title: string; description: string; type: 'newbie' | 'daily' | 'growth'; rewardPoints: number; targetValue: number; currentValue: number; status: 'todo' | 'done' | 'claimed'; actionText: string; actionUrl: string; }
export interface RankingItem { rank: number; targetId: number; targetType: 'user' | 'circle'; name: string; avatar: string; level?: number; levelName?: string; points?: number; postCount?: number; likeReceivedCount?: number; memberCount?: number; score: number; isCurrentUser?: boolean; }

export interface Activity { activityId: number; type: 'topic' | 'vote' | 'essay' | 'checkin'; title: string; description: string; coverImage: string; topicId?: number; topicName?: string; startAt: string; endAt: string; participantCount: number; submissionCount: number; status: 'not_started' | 'ongoing' | 'ended' | 'offline'; voteConfig?: { question: string; multiple: boolean; voted: boolean; selectedOptionIds: number[]; options: { optionId: number; text: string; voteCount: number; percent: number; }[] }; checkInConfig?: { totalDays: number; currentDays: number; todayChecked: boolean; todayTask: string; }; }
export interface OfficialAnnouncement { announcementId: number; title: string; summary: string; content: string; publisherName: string; status: 'draft' | 'published' | 'offline'; publishedAt?: string; createdAt: string; updatedAt: string; }
export interface Banner { bannerId: number; title: string; imageUrl: string; targetUrl: string; position: 'home' | 'hot' | 'circle' | 'topic' | 'activity'; weight: number; status: EnabledStatus; createdAt: string; }
export interface NotificationItem { notificationId: number; title: string; content: string; category: 'all' | 'interaction' | 'follow' | 'circle' | 'review' | 'growth' | 'activity' | 'system'; targetUrl?: string; readStatus: 'unread' | 'read'; createdAt: string; }

export interface BaseModel { modelId: number; provider: AiProvider; providerName: string; modelName: string; baseUrl: string; isDefault: boolean; status: EnabledStatus; createdAt: string; updatedAt: string; }
export interface TokenConfig { configId: number; dailyMaxTokens: number; userDailyMaxTokens: number; singleConversationMaxTokens: number; overLimitStrategy: 'reject' | 'downgrade' | 'warn'; updatedAt: string; }
export interface UserTokenUsage { userId: number; nickname: string; avatar: string; todayTokens: number; monthTokens: number; lastUsedAt: string; }
export interface AdminLog { logId: number; adminName: string; action: string; targetType: string; targetId: number; detail: string; createdAt: string; }

export interface HotPostRankItem { rank: number; postId: number; title: string; authorName: string; hotScore: number; likeCount: number; commentCount: number; }
export interface HotCircleRankItem { rank: number; circleId: number; name: string; memberCount: number; postCount: number; featuredPostCount: number; hotScore: number; }
export interface HotTopicRankItem { rank: number; topicId: number; name: string; participantCount: number; postCount: number; hotScore: number; }
