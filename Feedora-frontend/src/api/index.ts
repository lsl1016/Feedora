import type { Activity, AdminLog, AiChatMessage, AiChatSession, Banner, BaseModel, Comment, LoginResult, Circle, CircleAnnouncement, CircleChatMessage, CircleJoinRequest, CircleMember, CircleMemberRole, CircleSummarySubscription, ContentTag, CreatePostRequest, FollowingFeedItem, KnowledgeBase, KnowledgeBaseMember, MyCommentItem, NotificationItem, OfficialAnnouncement, PageRequest, PageResult, Post, PostQuery, RankingItem, Task, TokenConfig, Topic, User, UserModelConfig, UserTokenUsage, WorkspaceDashboardStats, WorkspaceNote } from '../types';
import { maskApiKey } from '../utils';
import { createApiAdapter } from './adapter';
import { httpDelete, httpGet, httpPost, httpPut } from './request';

const img = {
  ai: 'https://images.unsplash.com/photo-1677442136019-21780ecad995?auto=format&fit=crop&w=1200&q=80',
  code: 'https://images.unsplash.com/photo-1515879218367-8466d910aaa4?auto=format&fit=crop&w=1200&q=80',
  server: 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&w=1200&q=80',
  desk: 'https://images.unsplash.com/photo-1497366754035-f200968a6e72?auto=format&fit=crop&w=1200&q=80',
  front: 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?auto=format&fit=crop&w=1200&q=80',
};

function delay<T>(data: T, ms = 200): Promise<T> { return new Promise((resolve) => setTimeout(() => resolve(JSON.parse(JSON.stringify(data))), ms)); }
function rawDelay<T>(data: T, ms = 200): Promise<T> { return new Promise((resolve) => setTimeout(() => resolve(data), ms)); }
function now() { return '2026-07-05 18:30:00'; }
function page<T>(list: T[], pageNo = 1, pageSize = 10): PageResult<T> { return { list: list.slice((pageNo - 1) * pageSize, pageNo * pageSize), total: list.length, page: pageNo, pageSize }; }
function currentUserId() { try { return JSON.parse(localStorage.getItem('community_v21_user') || '{}').userId || 10001; } catch { return 10001; } }
function summary(u: User) { return { userId: u.userId, nickname: u.nickname, avatar: u.avatar, bio: u.bio, level: u.level, levelName: u.levelName }; }

export let users: User[] = [
  { userId: 10001, account: 'zhangsan', role: 'user', nickname: '张三', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan', bio: '后端开发，喜欢 Go、Java 和社区产品设计。', status: 'normal', level: 3, levelName: '内容贡献者', points: 560, experience: 720, nextLevelExperience: 800, badgeCount: 3, checkedInToday: false, continuousCheckInDays: 2, postCount: 12, commentCount: 35, followerCount: 256, followingCount: 86, likeReceivedCount: 1234, createdAt: '2026-07-01 10:00:00', updatedAt: now() },
  { userId: 10002, account: 'frontend_dev', role: 'user', nickname: '前端小李', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=frontend', bio: 'React / Vue 前端开发者，喜欢做清爽的页面。', status: 'normal', level: 4, levelName: '社区达人', points: 980, experience: 1200, nextLevelExperience: 2000, badgeCount: 5, checkedInToday: true, continuousCheckInDays: 7, postCount: 18, commentCount: 66, followerCount: 512, followingCount: 120, likeReceivedCount: 3450, createdAt: '2026-07-02 11:20:00', updatedAt: now() },
  { userId: 10003, account: 'admin', role: 'admin', nickname: '社区管理员', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin', bio: '社区运营与内容治理。', status: 'normal', level: 5, levelName: '核心创作者', points: 2400, experience: 2600, nextLevelExperience: 3000, badgeCount: 8, checkedInToday: true, continuousCheckInDays: 15, postCount: 6, commentCount: 12, followerCount: 900, followingCount: 30, likeReceivedCount: 5200, createdAt: '2026-06-30 09:00:00', updatedAt: now() },
  { userId: 10004, account: 'pm_wang', role: 'user', nickname: '产品经理小王', avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=pmwang', bio: '关注增长、社区、内容产品和用户体验。', status: 'normal', level: 3, levelName: '内容贡献者', points: 640, experience: 760, nextLevelExperience: 800, badgeCount: 4, checkedInToday: false, continuousCheckInDays: 1, postCount: 5, commentCount: 18, followerCount: 88, followingCount: 52, likeReceivedCount: 360, createdAt: '2026-07-03 08:00:00', updatedAt: now() },
];

export let tags: ContentTag[] = [
  { tagId: 1, tagName: 'Go语言', description: 'Go 后端开发、微服务、工程化实践。', status: 'enabled', useCount: 128, createdAt: '2026-07-01 10:00:00', updatedAt: now() },
  { tagId: 2, tagName: 'React', description: 'React 前端开发、组件设计、状态管理。', status: 'enabled', useCount: 96, createdAt: '2026-07-01 10:00:00', updatedAt: now() },
  { tagId: 3, tagName: 'AI工具', description: 'AI 编程助手、智能体、效率工具。', status: 'enabled', useCount: 180, createdAt: '2026-07-01 10:00:00', updatedAt: now() },
  { tagId: 4, tagName: '项目实战', description: '个人项目、MVP、工程实践。', status: 'enabled', useCount: 220, createdAt: '2026-07-01 10:00:00', updatedAt: now() },
  { tagId: 5, tagName: '面试求职', description: '简历、面试、求职经验。', status: 'enabled', useCount: 86, createdAt: '2026-07-01 10:00:00', updatedAt: now() },
];

export let topics: Topic[] = [
  { topicId: 1, name: 'AI编程助手实践', description: '分享 AI 编程助手在真实项目中的使用经验。', coverImage: img.ai, postCount: 560, participantCount: 2345, isOfficial: true, isRecommended: true, status: 'enabled', createdAt: '2026-07-01 10:00:00', updatedAt: now() },
  { topicId: 2, name: '个人项目部署', description: '讨论 Docker、Nginx、CI/CD、云服务器部署经验。', coverImage: img.server, postCount: 320, participantCount: 1200, isOfficial: true, isRecommended: true, status: 'enabled', createdAt: '2026-07-01 10:00:00', updatedAt: now() },
  { topicId: 3, name: 'React组件设计', description: '讨论组件抽象、状态管理和前端工程体验。', coverImage: img.front, postCount: 210, participantCount: 880, isOfficial: false, isRecommended: true, status: 'enabled', createdAt: '2026-07-02 10:00:00', updatedAt: now() },
];

export let circles: Circle[] = [
  { circleId: 1, name: 'Go 后端开发圈', avatar: 'https://api.dicebear.com/7.x/shapes/svg?seed=go-circle', description: '专注 Go、微服务、性能优化和工程实践。', category: '后端开发', tags: [{ tagId: 1, tagName: 'Go语言' }, { tagId: 4, tagName: '项目实战' }], ownerId: 10001, owner: summary(users[0]), joinType: 'direct', postPermission: 'all', memberCount: 12345, postCount: 3210, featuredPostCount: 120, isJoined: true, myRole: 'owner', myStatus: 'normal', isRecommended: true, status: 'normal', rules: '1. 禁止广告引流。\n2. 禁止人身攻击。\n3. 提问前请先搜索历史内容。', createdAt: '2026-07-01 10:00:00', updatedAt: now() },
  { circleId: 2, name: 'React 前端圈', avatar: 'https://api.dicebear.com/7.x/shapes/svg?seed=react-circle', description: '讨论 React、TypeScript、组件化和前端工程化。', category: '前端开发', tags: [{ tagId: 2, tagName: 'React' }, { tagId: 4, tagName: '项目实战' }], ownerId: 10002, owner: summary(users[1]), joinType: 'approval', postPermission: 'all', memberCount: 8800, postCount: 2100, featuredPostCount: 88, isJoined: false, isRecommended: true, status: 'normal', rules: '保持友好讨论，禁止无意义刷屏。', createdAt: '2026-07-01 10:00:00', updatedAt: now() },
  { circleId: 3, name: 'AI 工具分享圈', avatar: 'https://api.dicebear.com/7.x/shapes/svg?seed=ai-circle', description: '分享 AI 工具、智能体和效率实践。', category: 'AI', tags: [{ tagId: 3, tagName: 'AI工具' }], ownerId: 10003, owner: summary(users[2]), joinType: 'direct', postPermission: 'all', memberCount: 15600, postCount: 4100, featuredPostCount: 220, isJoined: true, myRole: 'moderator', myStatus: 'normal', isRecommended: true, status: 'normal', rules: '欢迎分享真实使用体验，禁止夸大宣传。', createdAt: '2026-07-02 10:00:00', updatedAt: now() },
];

export let posts: Post[] = [
  { postId: 20001, postType: 'original', authorId: 10002, author: summary(users[1]), title: '如何高效学习一门新的编程语言？', content: '最近在学习 Go 语言，感觉语法并不难，难点在于工程实践。建议先快速过一遍语法，然后直接做一个小项目，在项目中遇到问题再回头补知识点。', summary: '最近在学习 Go 语言，感觉语法并不难，难点在于工程实践...', images: [img.code, img.server], tags: [{ tagId: 1, tagName: 'Go语言' }, { tagId: 4, tagName: '项目实战' }], topics: [{ topicId: 2, name: '个人项目部署' }], circleId: 1, circle: { circleId: 1, name: 'Go 后端开发圈', avatar: circles[0].avatar, memberCount: circles[0].memberCount, postCount: circles[0].postCount }, visibility: 'public', status: 'published', isTop: false, isFeatured: true, isSelected: false, viewCount: 1280, likeCount: 432, commentCount: 32, favoriteCount: 88, shareCount: 24, repostCount: 12, hotScore: 3200, liked: false, favorited: false, followedAuthor: true, createdAt: '2026-07-05 09:30:00', publishedAt: '2026-07-05 09:35:00', updatedAt: now() },
  { postId: 20002, postType: 'original', authorId: 10001, author: summary(users[0]), title: '三层架构风格的社区前端设计分享', content: '面向开发者社区，前台采用内容流、左侧模块导航、右侧信息栏，后台采用管理系统结构。', summary: '面向开发者社区，前台采用内容流、左侧模块导航、右侧信息栏...', images: [img.desk], tags: [{ tagId: 2, tagName: 'React' }, { tagId: 4, tagName: '项目实战' }], topics: [{ topicId: 3, name: 'React组件设计' }], circleId: 2, circle: { circleId: 2, name: 'React 前端圈', avatar: circles[1].avatar, memberCount: circles[1].memberCount, postCount: circles[1].postCount }, visibility: 'public', status: 'published', isTop: true, isFeatured: false, isSelected: true, viewCount: 980, likeCount: 156, commentCount: 24, favoriteCount: 45, shareCount: 12, repostCount: 6, hotScore: 1380, liked: true, favorited: false, followedAuthor: false, createdAt: '2026-07-04 18:00:00', publishedAt: '2026-07-04 18:05:00', updatedAt: now() },
  { postId: 20003, postType: 'original', authorId: 10004, author: summary(users[3]), title: '做程序员社区时，冷启动内容应该怎么准备？', content: '社区冷启动最怕空场，需要准备高质量种子内容，覆盖典型问题、经验分享和讨论话题。', summary: '社区冷启动最怕空场，需要准备高质量种子内容...', images: [img.ai], tags: [{ tagId: 3, tagName: 'AI工具' }, { tagId: 4, tagName: '项目实战' }], topics: [{ topicId: 1, name: 'AI编程助手实践' }], circleId: 3, circle: { circleId: 3, name: 'AI 工具分享圈', avatar: circles[2].avatar, memberCount: circles[2].memberCount, postCount: circles[2].postCount }, visibility: 'public', status: 'published', isTop: false, isFeatured: true, isSelected: true, viewCount: 720, likeCount: 120, commentCount: 18, favoriteCount: 39, shareCount: 18, repostCount: 4, hotScore: 1020, liked: false, favorited: true, followedAuthor: false, createdAt: '2026-07-03 15:00:00', publishedAt: '2026-07-03 15:05:00', updatedAt: now() },
  { postId: 20004, postType: 'original', authorId: 10001, author: summary(users[0]), title: 'React 状态管理方案对比', content: '本文简单对比 Zustand、Redux Toolkit 和 Context 在中小型项目中的使用体验。', summary: '本文简单对比 Zustand、Redux Toolkit 和 Context...', images: [], tags: [{ tagId: 2, tagName: 'React' }], topics: [{ topicId: 3, name: 'React组件设计' }], circleId: 2, circle: { circleId: 2, name: 'React 前端圈', avatar: circles[1].avatar, memberCount: circles[1].memberCount, postCount: circles[1].postCount }, visibility: 'circle_only', status: 'hidden', isTop: false, isFeatured: false, isSelected: false, viewCount: 0, likeCount: 0, commentCount: 0, favoriteCount: 0, shareCount: 0, repostCount: 0, hotScore: 0, liked: false, favorited: false, followedAuthor: false, createdAt: '2026-07-05 14:00:00', updatedAt: now() },
];
posts.push({ ...posts[0], postId: 20005, postType: 'repost', authorId: 10001, author: summary(users[0]), title: `转发：${posts[0].title}`, content: '', summary: '这个学习方法很适合新手，先做项目再补知识点。', images: [], sourcePostId: 20001, sourcePost: posts[0], repostComment: '这个学习方法很适合新手，先做项目再补知识点。', likeCount: 80, commentCount: 12, favoriteCount: 10, shareCount: 5, repostCount: 2, hotScore: 600, createdAt: '2026-07-05 11:00:00', publishedAt: '2026-07-05 11:00:00' });

export let comments: Comment[] = [
  { commentId: 30001, postId: 20001, userId: 10001, user: summary(users[0]), content: '这个学习方法很实用，先做项目再补知识点效率更高。', likeCount: 12, liked: false, status: 'normal', replies: [], createdAt: '2026-07-05 10:30:00', updatedAt: now() },
  { commentId: 30002, postId: 20003, userId: 10001, user: summary(users[0]), content: '种子内容确实很关键，面向开发者的社区内容质量比数量更重要。', likeCount: 8, liked: false, status: 'normal', replies: [], createdAt: '2026-07-05 11:20:00', updatedAt: now() },
];

export let circleMembers: CircleMember[] = [
  { id: 1, circleId: 1, userId: 10001, user: summary(users[0]), role: 'owner', status: 'normal', joinedAt: '2026-07-01 10:00:00', updatedAt: now() },
  { id: 2, circleId: 1, userId: 10002, user: summary(users[1]), role: 'moderator', status: 'normal', joinedAt: '2026-07-02 10:00:00', updatedAt: now() },
  { id: 3, circleId: 1, userId: 10004, user: summary(users[3]), role: 'member', status: 'muted', muteReason: '刷屏', mutedUntil: '7 天后', joinedAt: '2026-07-03 10:00:00', updatedAt: now() },
];
export let joinRequests: CircleJoinRequest[] = [{ requestId: 1, circleId: 2, userId: 10001, user: summary(users[0]), reason: '想学习 React 工程化实践。', status: 'pending', createdAt: now() }];
export let announcements: CircleAnnouncement[] = [{ announcementId: 1, circleId: 1, title: '新人必读：圈子提问规范', content: '提问前请先搜索历史内容，描述问题时请包含环境、现象和复现步骤。', publisher: summary(users[0]), createdAt: now() }];
export let chatMessages: CircleChatMessage[] = [
  { messageId: 1, circleId: 1, senderId: 0, sender: { userId: 0, nickname: '系统', avatar: '' }, content: '欢迎加入 Go 后端开发圈，请遵守圈子规则。', messageType: 'system', status: 'normal', createdAt: '2026-07-05 09:00:00' },
  { messageId: 2, circleId: 1, senderId: 10001, sender: summary(users[0]), content: '今天有人看过 React Server Components 吗？', messageType: 'text', status: 'normal', createdAt: '2026-07-05 10:30:00' },
  { messageId: 3, circleId: 1, senderId: 10002, sender: summary(users[1]), content: '看过，建议先理解服务端渲染边界。', messageType: 'text', status: 'normal', createdAt: '2026-07-05 10:31:00' },
];
export let summarySubscriptions: CircleSummarySubscription[] = [
  { subscriptionId: 1, userId: 10001, circleId: 1, circleName: 'Go 后端开发圈', enabled: true, frequency: 'daily', pushTime: '09:00', channels: ['notification', 'email'], email: 'user@example.com', contentScopes: ['post', 'comment', 'chat', 'featured'], createdAt: now(), updatedAt: now() },
  { subscriptionId: 2, userId: 10001, circleId: 3, circleName: 'AI 工具分享圈', enabled: true, frequency: 'once', oneTimePushAt: '2026-07-05 18:00:00', channels: ['notification'], contentScopes: ['post', 'chat', 'featured'], createdAt: now(), updatedAt: now() },
];

export let notes: WorkspaceNote[] = [
  { noteId: 1, title: 'React 状态管理学习笔记', content: '整理 Zustand、Redux Toolkit、Context 的适用场景。', summary: '整理 Zustand、Redux Toolkit、Context 的适用场景。', knowledgeBaseId: 1, knowledgeBaseName: '前端工程化知识库', ownerId: 10001, status: 'normal', createdAt: '2026-07-05 10:00:00', updatedAt: '2026-07-05 18:30:00' },
  { noteId: 2, title: 'Go 项目 Docker 部署记录', content: '记录 Go 项目使用 Docker Compose、Nginx、MySQL、Redis 部署的流程。', summary: '记录 Go 项目 Docker Compose 部署流程。', knowledgeBaseId: 2, knowledgeBaseName: 'Go 后端工程化知识库', ownerId: 10001, status: 'normal', createdAt: '2026-07-04 10:00:00', updatedAt: '2026-07-05 11:30:00' },
];
export let knowledgeBases: KnowledgeBase[] = [
  { knowledgeBaseId: 1, name: '前端工程化知识库', description: '记录 React、Vue、TypeScript、组件设计和工程化实践。', ownerId: 10001, noteCount: 8, memberCount: 2, visibility: 'collaborative', createdAt: '2026-07-01 10:00:00', updatedAt: now() },
  { knowledgeBaseId: 2, name: 'Go 后端工程化知识库', description: '记录 Go 项目实战、微服务、部署和性能优化经验。', ownerId: 10001, noteCount: 12, memberCount: 3, visibility: 'collaborative', createdAt: '2026-07-02 10:00:00', updatedAt: now() },
];
export let kbMembers: KnowledgeBaseMember[] = [
  { id: 1, knowledgeBaseId: 2, userId: 10001, user: summary(users[0]), role: 'owner', status: 'active', joinedAt: '2026-07-02 10:00:00' },
  { id: 2, knowledgeBaseId: 2, userId: 10002, user: summary(users[1]), role: 'editor', status: 'active', joinedAt: '2026-07-03 10:00:00' },
];
export let modelConfig: UserModelConfig = { configId: 1, userId: 10001, provider: 'deepseek', baseUrl: 'https://api.deepseek.com', apiKeyMasked: 'sk-****abcd', modelName: 'deepseek-chat', temperature: 0.7, maxTokens: 4096, status: 'enabled', createdAt: now(), updatedAt: now() };
export let aiSessions: AiChatSession[] = [{ sessionId: 1, title: 'Java 后端面试准备路线', userId: 10001, createdAt: now(), updatedAt: now() }];
export let aiMessages: AiChatMessage[] = [{ messageId: 1, sessionId: 1, role: 'user', content: '如何准备 Java 后端面试？', editable: false, createdAt: now() }, { messageId: 2, sessionId: 1, role: 'assistant', content: '可以从 Java 基础、JVM、并发、Spring、数据库、Redis、项目经验和系统设计几个方向准备。', editable: true, createdAt: now() }];

export let tasks: Task[] = [{ taskId: 1, title: '完善个人资料', description: '上传头像并填写个人简介。', type: 'newbie', rewardPoints: 20, targetValue: 1, currentValue: 0, status: 'todo', actionText: '去完成', actionUrl: '/users/me' }, { taskId: 2, title: '发布第一篇帖子', description: '完成一次内容创作。', type: 'newbie', rewardPoints: 30, targetValue: 1, currentValue: 1, status: 'done', actionText: '领取奖励', actionUrl: '/posts/create' }, { taskId: 3, title: '每日签到', description: '每天签到获得积分。', type: 'daily', rewardPoints: 5, targetValue: 1, currentValue: 1, status: 'done', actionText: '领取奖励', actionUrl: '/tasks' }];
export let badges = [{ badgeId: 1, name: '首次发帖', icon: '🥇', description: '发布第一篇帖子后获得。', condition: '发布第一篇帖子', status: 'enabled', createdAt: now() }, { badgeId: 2, name: '人气作者', icon: '🔥', description: '累计获得 1000 个点赞。', condition: '累计获得 1000 个点赞', status: 'enabled', createdAt: now() }];
export let activities: Activity[] = [{ activityId: 1, type: 'topic', title: '本周讨论主题：AI 编程助手实践', description: '分享 AI 编程助手在真实项目中的使用经验。', coverImage: img.ai, topicId: 1, topicName: 'AI编程助手实践', startAt: '2026-07-01 00:00:00', endAt: '2026-07-15 23:59:59', participantCount: 2345, submissionCount: 560, status: 'ongoing' }, { activityId: 2, type: 'vote', title: '你最常用的前端框架是？', description: '选择你当前最常用的前端框架。', coverImage: img.front, startAt: '2026-07-01 00:00:00', endAt: '2026-07-10 23:59:59', participantCount: 1024, submissionCount: 1024, status: 'ongoing', voteConfig: { question: '你最常用的前端框架是？', multiple: false, voted: false, selectedOptionIds: [], options: [{ optionId: 1, text: 'React', voteCount: 560, percent: 55 }, { optionId: 2, text: 'Vue', voteCount: 388, percent: 38 }] } }];
export let officialAnnouncements: OfficialAnnouncement[] = [{ announcementId: 1, title: '社区 V2.1 功能上线说明', summary: '工作空间、AI 问答、圈子聊天和智能摘要已上线。', content: '本次更新上线工作空间、笔记、知识库、AI 问答、圈子聊天、圈子智能摘要和后台 AI 配置能力。欢迎体验。', publisherName: '社区管理员', status: 'published', publishedAt: now(), createdAt: now(), updatedAt: now() }];
export let banners: Banner[] = [{ bannerId: 1, title: '开发者知识社区：AI + 知识库 + 圈子协作', imageUrl: img.ai, targetUrl: '/activities/1', position: 'home', weight: 100, status: 'enabled', createdAt: now() }];
export let notifications: NotificationItem[] = [{ notificationId: 1, title: '圈子加入申请已通过', content: '你的申请已通过，欢迎加入 Go 后端开发圈。', category: 'circle', targetUrl: '/circles/1', readStatus: 'unread', createdAt: now() }, { notificationId: 2, title: '任务奖励待领取', content: '发布第一篇帖子任务已完成，快去领取奖励。', category: 'growth', targetUrl: '/tasks', readStatus: 'unread', createdAt: now() }];
export let baseModels: BaseModel[] = [{ modelId: 1, provider: 'deepseek', providerName: 'DeepSeek', modelName: 'deepseek-chat', baseUrl: 'https://api.deepseek.com', isDefault: true, status: 'enabled', createdAt: now(), updatedAt: now() }, { modelId: 2, provider: 'openai', providerName: 'OpenAI', modelName: 'gpt-4o-mini', baseUrl: 'https://api.openai.com/v1', isDefault: false, status: 'disabled', createdAt: now(), updatedAt: now() }];
export let tokenConfig: TokenConfig = { configId: 1, dailyMaxTokens: 1000000, userDailyMaxTokens: 50000, singleConversationMaxTokens: 4096, overLimitStrategy: 'reject', updatedAt: now() };
export let tokenUsages: UserTokenUsage[] = [{ userId: 10001, nickname: '张三', avatar: users[0].avatar, todayTokens: 12800, monthTokens: 168000, lastUsedAt: now() }, { userId: 10002, nickname: '前端小李', avatar: users[1].avatar, todayTokens: 8600, monthTokens: 96000, lastUsedAt: now() }];
export let adminLogs: AdminLog[] = [{ logId: 1, adminName: '社区管理员', action: '更新 Token 配置', targetType: 'ai', targetId: 1, detail: '调整单用户每日 Token 限额', createdAt: now() }];

function aiAnswer(question: string) { if (question.includes('面试')) return '可以从基础知识、项目经验、系统设计、数据库、缓存和表达能力几个方向准备。建议把简历项目拆成背景、技术选型、难点、指标和结果来表达。'; if (question.includes('React')) return 'React 学习可以分为 JSX、组件、Hooks、状态管理、路由、工程化和项目实战几个阶段。'; if (question.includes('Go')) return 'Go 学习建议从语法、并发、Gin、Gorm、Redis、Docker 部署和微服务实践开始。'; return '建议先明确目标，再拆分学习路径，最后通过一个小项目进行验证和沉淀。'; }

export const mockApi = {
  delay,
  login: async (account: string, password: string): Promise<LoginResult> => { const user = users.find((u) => u.account === account) || (account === 'admin' ? users[2] : users[0]); if (!password) throw new Error('请输入密码'); return delay({ token: `token-${user.userId}`, user }); },
  register: async (data: { account: string; nickname: string; password: string; confirmPassword: string }) => { if (data.password !== data.confirmPassword) throw new Error('两次密码不一致'); const user: User = { ...users[0], userId: Math.max(...users.map((u) => u.userId)) + 1, account: data.account, nickname: data.nickname, avatar: `https://api.dicebear.com/7.x/avataaars/svg?seed=${data.account}`, role: 'user', points: 0, experience: 0, postCount: 0, followerCount: 0, likeReceivedCount: 0, createdAt: now(), updatedAt: now() }; users.unshift(user); return delay(user); },
  getUsers: (params: PageRequest & { keyword?: string }) => delay(page(users.filter((u) => !params.keyword || u.nickname.includes(params.keyword)), params.page, params.pageSize)),
  updateUser: async (user: User) => { const i = users.findIndex((u) => u.userId === user.userId); if (i >= 0) users[i] = user; return delay(user); },
  getTags: () => delay(tags),
  getTopics: (params: PageRequest & { type?: string }) => delay(page(topics.filter((t) => params.type === 'official' ? t.isOfficial : params.type === 'hot' ? true : true).sort((a,b)=>params.type==='hot'?b.participantCount-a.participantCount:0), params.page, params.pageSize)),
  getTopic: (id: number) => delay(topics.find((t) => t.topicId === id) || topics[0]),
  getCircles: (params: PageRequest & { scope?: string; keyword?: string; sort?: string }) => { let list = circles.filter((c) => c.status !== 'deleted'); if (params.scope === 'recommended') list = list.filter((c) => c.isRecommended); if (params.scope === 'joined') list = list.filter((c) => c.isJoined); if (params.scope === 'created') list = list.filter((c) => c.ownerId === currentUserId()); if (params.keyword) list = list.filter((c) => c.name.includes(params.keyword || '')); return delay(page(list, params.page, params.pageSize)); },
  getCircle: (id: number) => delay(circles.find((c) => c.circleId === id) || circles[0]),
  createCircle: (data: Partial<Circle>) => { const owner = users.find((u)=>u.userId===currentUserId()) || users[0]; const c: Circle = { circleId: Math.max(...circles.map((x)=>x.circleId))+1, name: data.name || '新圈子', avatar: data.avatar || 'https://api.dicebear.com/7.x/shapes/svg?seed=new-circle', description: data.description || '', category: data.category || '开发', tags: tags.slice(0,2).map((t)=>({tagId:t.tagId,tagName:t.tagName})), ownerId: owner.userId, owner: summary(owner), joinType: data.joinType || 'direct', postPermission: data.postPermission || 'all', memberCount: 1, postCount: 0, featuredPostCount: 0, isJoined: true, myRole: 'owner', myStatus: 'normal', isRecommended: false, status: 'normal', rules: data.rules || '', createdAt: now(), updatedAt: now() }; circles.unshift(c); return delay(c); },
  joinCircle: (circleId: number, reason?: string) => { const c = circles.find((x)=>x.circleId===circleId); if (!c) throw new Error('圈子不存在'); if (c.joinType === 'approval') { joinRequests.unshift({ requestId: Math.max(0,...joinRequests.map(r=>r.requestId))+1, circleId, userId: currentUserId(), user: summary(users.find(u=>u.userId===currentUserId()) || users[0]), reason: reason || '想加入圈子', status: 'pending', createdAt: now() }); return delay({ pending: true }); } c.isJoined = true; c.memberCount++; return delay({ pending: false }); },
  getPosts: (params: PostQuery) => { let list = posts; if (!params.includeHidden) list = list.filter((p) => p.status === 'published' && p.visibility === 'public'); if (params.includeHidden) list = list.filter((p) => params.status === 'all' || !params.status || p.status === params.status); if (params.feedType === 'following') list = list.filter((p) => p.followedAuthor || p.circleId === 1 || p.circleId === 3); if (params.tagId) list = list.filter((p) => p.tags.some((t)=>t.tagId===params.tagId)); if (params.topicId) list = list.filter((p) => p.topics.some((t)=>t.topicId===params.topicId)); if (params.circleId) list = list.filter((p) => p.circleId === params.circleId); if (params.keyword) list = list.filter((p) => `${p.title}${p.summary}`.includes(params.keyword || '')); list = [...list].sort((a,b)=> (params.sort==='hot'||params.feedType==='hot') ? b.hotScore-a.hotScore : Date.parse(b.createdAt)-Date.parse(a.createdAt)); return delay(page(list, params.page, params.pageSize)); },
  getHotRanks: (rankType: 'post' | 'circle' | 'topic', timeRange: 'today' | 'week' | 'all') => {
    const boost = timeRange === 'today' ? 1 : timeRange === 'week' ? 1.15 : 1.35;
    if (rankType === 'circle') return delay(circles.map((c, i)=>({ rank:i+1, circleId:c.circleId, name:c.name, memberCount:c.memberCount, postCount:c.postCount, featuredPostCount:c.featuredPostCount, hotScore:Math.round((c.memberCount*0.4+c.postCount*1.2+c.featuredPostCount*8)*boost) })).sort((a,b)=>b.hotScore-a.hotScore).map((x,i)=>({...x,rank:i+1})));
    if (rankType === 'topic') return delay(topics.map((t, i)=>({ rank:i+1, topicId:t.topicId, name:t.name, participantCount:t.participantCount, postCount:t.postCount, hotScore:Math.round((t.participantCount*1.2+t.postCount*6)*boost) })).sort((a,b)=>b.hotScore-a.hotScore).map((x,i)=>({...x,rank:i+1})));
    return delay(posts.filter(p=>p.status==='published').map((p, i)=>({ rank:i+1, postId:p.postId, title:p.title, authorName:p.author.nickname, hotScore:Math.round(p.hotScore*boost), likeCount:p.likeCount, commentCount:p.commentCount })).sort((a,b)=>b.hotScore-a.hotScore).map((x,i)=>({...x,rank:i+1})));
  },
  getPost: (id: number) => delay(posts.find((p)=>p.postId===id) || posts[0]),
  createPost: (data: CreatePostRequest) => { const u = users.find((x)=>x.userId===currentUserId()) || users[0]; const status = data.publishMode === 'draft' ? 'draft' : data.publishMode === 'schedule' ? 'scheduled' : 'published'; const c = data.circleId ? circles.find((x)=>x.circleId===data.circleId) : undefined; const p: Post = { postId: Math.max(...posts.map((x)=>x.postId))+1, postType: 'original', authorId: u.userId, author: summary(u), title: data.title, content: data.content, summary: data.content.slice(0,80), images: data.images, tags: tags.filter((t)=>data.tagIds.includes(t.tagId)).map((t)=>({tagId:t.tagId, tagName:t.tagName})), topics: topics.filter((t)=>data.topicIds.includes(t.topicId)).map((t)=>({topicId:t.topicId, name:t.name})), circleId: data.circleId, circle: c ? { circleId: c.circleId, name: c.name, avatar: c.avatar, memberCount: c.memberCount, postCount: c.postCount } : undefined, visibility: data.visibility, status, isTop: false, isFeatured: false, isSelected: false, scheduledAt: data.scheduledAt, viewCount: 0, likeCount: 0, commentCount: 0, favoriteCount: 0, shareCount: 0, repostCount: 0, hotScore: 0, liked: false, favorited: false, followedAuthor: false, createdAt: now(), publishedAt: status==='published'?now():undefined, updatedAt: now() }; posts.unshift(p); return delay(p); },
  updatePost: (id: number, patch: Partial<Post>) => { const p = posts.find((x)=>x.postId===id); if (p) Object.assign(p, patch, { updatedAt: now() }); return delay(p); },
  hidePost: (id: number) => { const p = posts.find((x)=>x.postId===id); if (p) p.status='hidden'; return delay({ postId: id, status: 'hidden' as const }); },
  unhidePost: (id: number) => { const p = posts.find((x)=>x.postId===id); if (p) p.status='published'; return delay({ postId: id, status: 'published' as const }); },
  deletePost: (id: number) => { const p = posts.find((x)=>x.postId===id); if (p) p.status='deleted'; return delay(true); },
  likePost: (id: number) => { const p = posts.find((x)=>x.postId===id); if (p) { p.liked=!p.liked; p.likeCount += p.liked ? 1 : -1; } return delay(p); },
  favoritePost: (id: number) => { const p = posts.find((x)=>x.postId===id); if (p) { p.favorited=!p.favorited; p.favoriteCount += p.favorited ? 1 : -1; } return delay(p); },
  sharePost: (id: number) => { const p = posts.find((x)=>x.postId===id); if (p) p.shareCount++; return delay(p); },
  repostPost: (sourcePostId: number, repostComment: string) => { const source = posts.find((p)=>p.postId===sourcePostId) || posts[0]; const u = users.find((x)=>x.userId===currentUserId()) || users[0]; const p: Post = { ...source, postId: Math.max(...posts.map((x)=>x.postId))+1, postType: 'repost', authorId: u.userId, author: summary(u), title: `转发：${source.title}`, content: '', summary: repostComment || source.summary, images: [], sourcePostId, sourcePost: source, repostComment, likeCount:0, commentCount:0, favoriteCount:0, shareCount:0, repostCount:0, hotScore:0, liked:false, favorited:false, createdAt:now(), publishedAt:now() }; posts.unshift(p); source.repostCount++; return delay(p); },
  getComments: (postId: number) => delay(page(comments.filter((c)=>c.postId===postId && c.status==='normal'),1,50)),
  addComment: (postId: number, content: string) => { const u=users.find((x)=>x.userId===currentUserId())||users[0]; const c: Comment={commentId:Math.max(0,...comments.map(x=>x.commentId))+1,postId,userId:u.userId,user:summary(u),content,likeCount:0,liked:false,status:'normal',replies:[],createdAt:now(),updatedAt:now()}; comments.unshift(c); const p=posts.find(x=>x.postId===postId); if(p)p.commentCount++; return delay(c); },
  getMyComments: (params: PageRequest) => delay(page(comments.filter((c)=>c.userId===currentUserId()).map((c)=>({commentId:c.commentId, postId:c.postId, postTitle:posts.find((p)=>p.postId===c.postId)?.title||'原帖已删除', content:c.content, likeCount:c.likeCount, status:c.status, createdAt:c.createdAt})), params.page, params.pageSize)),
  deleteComment: (commentId: number) => { const c=comments.find(x=>x.commentId===commentId); if(c)c.status='deleted'; return delay(true); },
  getCircleMembers: (circleId: number) => delay(page(circleMembers.filter((m)=>m.circleId===circleId && m.status!=='removed'),1,100)),
  setMemberRole: (circleId:number,userId:number,role:CircleMemberRole)=>{const m=circleMembers.find(x=>x.circleId===circleId&&x.userId===userId); if(m)m.role=role; return delay(true);},
  muteMember: (circleId:number,userId:number,duration:string,reason:string)=>{const m=circleMembers.find(x=>x.circleId===circleId&&x.userId===userId); if(m){m.status='muted';m.muteReason=reason;m.mutedUntil=duration;} return delay(true);},
  unmuteMember: (circleId:number,userId:number)=>{const m=circleMembers.find(x=>x.circleId===circleId&&x.userId===userId); if(m)m.status='normal'; return delay(true);},
  removeMember: (circleId:number,userId:number)=>{const m=circleMembers.find(x=>x.circleId===circleId&&x.userId===userId); if(m)m.status='removed'; return delay(true);},
  getJoinRequests: (circleId:number)=>delay(joinRequests.filter((r)=>r.circleId===circleId&&r.status==='pending')),
  handleJoinRequest: (id:number,status:'approved'|'rejected')=>{const r=joinRequests.find(x=>x.requestId===id); if(r)r.status=status; return delay(true);},
  getAnnouncements: (circleId:number)=>delay(announcements.filter((a)=>a.circleId===circleId)),
  addAnnouncement: (circleId:number,title:string,content:string)=>{const a={announcementId:Math.max(0,...announcements.map(x=>x.announcementId))+1,circleId,title,content,publisher:summary(users[0]),createdAt:now()}; announcements.unshift(a); return delay(a);},
  getChatMessages: (circleId:number)=>delay(chatMessages.filter((m)=>m.circleId===circleId)),
  sendChatMessage: (circleId:number,content:string)=>{const u=users.find((x)=>x.userId===currentUserId())||users[0]; const m:CircleChatMessage={messageId:Math.max(0,...chatMessages.map(x=>x.messageId))+1,circleId,senderId:u.userId,sender:summary(u),content,messageType:'text',status:'normal',createdAt:now()}; chatMessages.push(m); return delay(m);},
  getSummarySubscriptions: ()=>delay(summarySubscriptions.filter((s)=>s.userId===currentUserId() && s.enabled)),
  getSummarySubscription: (circleId:number)=>delay(summarySubscriptions.find((s)=>s.circleId===circleId&&s.userId===currentUserId()) || { subscriptionId:0,userId:currentUserId(),circleId,circleName:circles.find(c=>c.circleId===circleId)?.name,enabled:false,frequency:'daily',pushTime:'09:00',channels:['notification'],contentScopes:['post','chat'],createdAt:now(),updatedAt:now() } as CircleSummarySubscription),
  saveSummarySubscription: (circleId:number,data:Partial<CircleSummarySubscription>)=>{let s=summarySubscriptions.find(x=>x.circleId===circleId&&x.userId===currentUserId()); if(!s){s={subscriptionId:Math.max(0,...summarySubscriptions.map(x=>x.subscriptionId))+1,userId:currentUserId(),circleId,circleName:circles.find(c=>c.circleId===circleId)?.name,enabled:false,frequency:'daily',pushTime:'09:00',channels:['notification'],contentScopes:['post'],createdAt:now(),updatedAt:now()}; summarySubscriptions.push(s);} Object.assign(s,data,{circleName:circles.find(c=>c.circleId===circleId)?.name,updatedAt:now()}); return delay(s);},
  cancelSummarySubscription: (circleId:number)=>{const s=summarySubscriptions.find(x=>x.circleId===circleId&&x.userId===currentUserId()); if(s)s.enabled=false; return delay(true);},
  getBanners: () => delay(banners.filter((b)=>b.status==='enabled')),
  getActivities: (params: PageRequest & { type?: string }) => delay(page(activities.filter((a)=>!params.type||params.type==='all'||a.type===params.type), params.page, params.pageSize)),
  getActivity: (id:number)=>delay(activities.find(a=>a.activityId===id)||activities[0]),
  getOfficialAnnouncements: (params:PageRequest)=>delay(page(officialAnnouncements.filter(a=>a.status==='published'),params.page,params.pageSize)),
  getNotificationList: ()=>delay(notifications),
  markNotificationRead: (id:number)=>{const n=notifications.find(x=>x.notificationId===id); if(n)n.readStatus='read'; return delay(true);},
  getFollowingUsers: (params:PageRequest)=>delay(page(users.filter((u)=>u.userId!==currentUserId()).slice(0,3),params.page,params.pageSize)),
  getFollowingCircles: (params:PageRequest)=>delay(page(circles.filter(c=>c.isJoined),params.page,params.pageSize)),
  getFollowingTopics: (params:PageRequest)=>delay(page(topics.slice(0,3),params.page,params.pageSize)),
  getFollowingTags: (params:PageRequest)=>delay(page(tags.slice(0,4),params.page,params.pageSize)),
  getFollowingFeed: (params:PageRequest & { feedTab?: string })=>delay(page([{feedId:1,feedType:'post',title:'前端小李 发布了帖子',summary:'如何高效学习一门新的编程语言？',targetId:20001,targetUrl:'/posts/20001',sourceName:'前端小李',sourceAvatar:users[1].avatar,createdAt:now()},{feedId:2,feedType:'circle',title:'Go 后端开发圈 有新精华内容',summary:'Go 项目部署复盘',targetId:1,targetUrl:'/circles/1',sourceName:'Go 后端开发圈',sourceAvatar:circles[0].avatar,createdAt:now()},{feedId:3,feedType:'topic',title:'#AI编程助手实践# 有新的热门讨论',summary:'AI 编程助手到底适合哪些场景？',targetId:1,targetUrl:'/topics/1',sourceName:'AI编程助手实践',createdAt:now()}] as FollowingFeedItem[],params.page,params.pageSize)),
  search: (keyword:string,type:string,params:PageRequest)=>{const all:any[]=[]; if(type==='all'||type==='post') all.push(...posts.filter(p=>p.status==='published'&&`${p.title}${p.summary}`.includes(keyword))); if(type==='all'||type==='user') all.push(...users.filter(u=>`${u.nickname}${u.bio}`.includes(keyword))); if(type==='all'||type==='topic') all.push(...topics.filter(t=>`${t.name}${t.description}`.includes(keyword))); if(type==='all'||type==='circle') all.push(...circles.filter(c=>`${c.name}${c.description}`.includes(keyword))); return delay(page(all,params.page,params.pageSize));},
  searchSuggest: (keyword:string)=>delay([...posts.filter(p=>p.title.includes(keyword)).slice(0,3).map(p=>({type:'post',id:p.postId,title:p.title,description:'帖子'})),...topics.filter(t=>t.name.includes(keyword)).slice(0,2).map(t=>({type:'topic',id:t.topicId,title:`#${t.name}#`,description:'话题'})),...circles.filter(c=>c.name.includes(keyword)).slice(0,2).map(c=>({type:'circle',id:c.circleId,title:c.name,description:'圈子'}))]),
  getHotKeywords:()=>delay(['AI 编程助手','Go 后端工程化','React 项目实战','Java 面试','个人项目部署']),
  getWorkspaceDashboard: ():Promise<WorkspaceDashboardStats> => delay({ noteCount: notes.filter(n=>n.status==='normal').length, knowledgeBaseCount: knowledgeBases.length, aiChatCount: aiSessions.length, collaborationMemberCount: kbMembers.length }),
  getNotes: (params:PageRequest & { keyword?:string; kbId?:number })=>delay(page(notes.filter(n=>n.status==='normal'&&(!params.keyword||n.title.includes(params.keyword))&&(!params.kbId||n.knowledgeBaseId===params.kbId)),params.page,params.pageSize)),
  getNote: (id:number)=>delay(notes.find(n=>n.noteId===id)||notes[0]),
  saveNote: (data:Partial<WorkspaceNote>)=>{ if(data.noteId){const n=notes.find(x=>x.noteId===data.noteId); if(n)Object.assign(n,data,{updatedAt:now()}); return delay(n!);} const n:WorkspaceNote={noteId:Math.max(0,...notes.map(x=>x.noteId))+1,title:data.title||'未命名笔记',content:data.content||'',summary:(data.content||'').slice(0,80),knowledgeBaseId:data.knowledgeBaseId,knowledgeBaseName:knowledgeBases.find(k=>k.knowledgeBaseId===data.knowledgeBaseId)?.name,ownerId:currentUserId(),status:'normal',createdAt:now(),updatedAt:now()}; notes.unshift(n); return delay(n); },
  deleteNote: (id:number)=>{const n=notes.find(x=>x.noteId===id); if(n)n.status='deleted'; return delay(true);},
  publishNoteAsPost: (id:number)=>delay({redirectUrl:`/posts/create?sourceType=note&noteId=${id}`}),
  optimizeNote: (content:string,goal:string)=>delay(`AI 优化结果（${goal}）：\n${aiAnswer(content)}\n\n核心要点：\n1. 先明确目标。\n2. 再拆解知识点。\n3. 最后通过项目实践沉淀。`),
  getKnowledgeBases: (params:PageRequest)=>delay(page(knowledgeBases,params.page,params.pageSize)),
  getKnowledgeBase: (id:number)=>delay(knowledgeBases.find(k=>k.knowledgeBaseId===id)||knowledgeBases[0]),
  saveKnowledgeBase: (data:Partial<KnowledgeBase>)=>{ if(data.knowledgeBaseId){const k=knowledgeBases.find(x=>x.knowledgeBaseId===data.knowledgeBaseId); if(k)Object.assign(k,data,{updatedAt:now()}); return delay(k!);} const k:KnowledgeBase={knowledgeBaseId:Math.max(0,...knowledgeBases.map(x=>x.knowledgeBaseId))+1,name:data.name||'新知识库',description:data.description||'',ownerId:currentUserId(),noteCount:0,memberCount:1,visibility:data.visibility||'private',createdAt:now(),updatedAt:now()}; knowledgeBases.unshift(k); return delay(k);},
  getKnowledgeBaseMembers:(kbId:number)=>delay(kbMembers.filter(m=>m.knowledgeBaseId===kbId)),
  inviteKbMember:(kbId:number,userId:number,role:'editor'|'viewer')=>{const u=users.find(x=>x.userId===userId)||users[1]; kbMembers.push({id:Math.max(0,...kbMembers.map(x=>x.id))+1,knowledgeBaseId:kbId,userId:userId,user:summary(u),role,status:'pending',invitedAt:now()}); return delay(true);},
  updateKbMemberRole:(kbId:number,userId:number,role:'owner'|'editor'|'viewer')=>{const m=kbMembers.find(x=>x.knowledgeBaseId===kbId&&x.userId===userId); if(m)m.role=role; return delay(true);},
  removeKbMember:(kbId:number,userId:number)=>{kbMembers=kbMembers.filter(m=>!(m.knowledgeBaseId===kbId&&m.userId===userId)); return delay(true);},
  getModelConfig:()=>delay(modelConfig),
  saveModelConfig:(data:Partial<UserModelConfig>&{apiKey?:string})=>{Object.assign(modelConfig,data,{apiKeyMasked:data.apiKey?maskApiKey(data.apiKey):modelConfig.apiKeyMasked,updatedAt:now()}); return delay(modelConfig);},
  testModelConfig:()=>delay(true),
  getAiSessions:()=>delay(aiSessions),
  createAiSession:(title='新会话')=>{const s:AiChatSession={sessionId:Math.max(0,...aiSessions.map(x=>x.sessionId))+1,title,userId:currentUserId(),createdAt:now(),updatedAt:now()}; aiSessions.unshift(s); return delay(s);},
  getAiMessages:(sessionId:number)=>delay(aiMessages.filter(m=>m.sessionId===sessionId)),
  sendAiMessage:(sessionId:number,content:string)=>{const userMsg:AiChatMessage={messageId:Math.max(0,...aiMessages.map(x=>x.messageId))+1,sessionId,role:'user',content,editable:false,createdAt:now()}; const answer:AiChatMessage={messageId:userMsg.messageId+1,sessionId,role:'assistant',content:aiAnswer(content),editable:true,createdAt:now()}; aiMessages.push(userMsg,answer); return delay([userMsg,answer]);},
  updateAiMessage:(id:number,content:string)=>{const m=aiMessages.find(x=>x.messageId===id); if(m)m.content=content; return delay(m);},
  saveAiToNote:(messageId:number,title:string,kbId?:number)=>{const m=aiMessages.find(x=>x.messageId===messageId); return mockApi.saveNote({title,content:m?.content||'',knowledgeBaseId:kbId});},
  topicAgent:(topicId:number,action:string,question?:string)=>delay({answer: action==='recommend_posts'?'推荐阅读：如何高效学习一门新的编程语言？、做社区产品时冷启动内容怎么准备？': aiAnswer(question||action), relatedPostIds:[20001,20003]}),
  checkIn:()=>{const u=users.find(x=>x.userId===currentUserId())||users[0]; if(!u.checkedInToday){u.checkedInToday=true;u.points+=5;u.continuousCheckInDays++;} return delay({points:5,continuousDays:u.continuousCheckInDays});},
  getTasks:(type:string)=>delay(tasks.filter(t=>t.type===type)),
  claimTask:(id:number)=>{const t=tasks.find(x=>x.taskId===id); if(t)t.status='claimed'; return delay(true);},
  getBadges:()=>delay(badges),
  getRankings:(type:string, range?:string)=>delay((type==='circle'?circles.map((c,i)=>({rank:i+1,targetId:c.circleId,targetType:'circle' as const,name:c.name,avatar:c.avatar,memberCount:c.memberCount,score:c.memberCount+c.postCount})):users.map((u,i)=>({rank:i+1,targetId:u.userId,targetType:'user' as const,name:u.nickname,avatar:u.avatar,level:u.level,levelName:u.levelName,points:u.points,postCount:u.postCount,likeReceivedCount:u.likeReceivedCount,score:u.points+u.likeReceivedCount,isCurrentUser:u.userId===currentUserId()}))).sort((a,b)=>b.score-a.score).map((x,i)=>({...x,rank:i+1})) as RankingItem[]),
  getBaseModels:(params:PageRequest)=>delay(page(baseModels,params.page,params.pageSize)),
  saveBaseModel:(data:Partial<BaseModel>)=>{ if(data.modelId){const m=baseModels.find(x=>x.modelId===data.modelId); if(m)Object.assign(m,data,{updatedAt:now()}); return delay(m!);} const m:BaseModel={modelId:Math.max(0,...baseModels.map(x=>x.modelId))+1,provider:data.provider||'custom',providerName:data.providerName||'自定义',modelName:data.modelName||'custom-model',baseUrl:data.baseUrl||'',isDefault:Boolean(data.isDefault),status:data.status||'enabled',createdAt:now(),updatedAt:now()}; if(m.isDefault) baseModels.forEach(x=>x.isDefault=false); baseModels.unshift(m); return delay(m);},
  setDefaultModel:(id:number)=>{baseModels.forEach(m=>m.isDefault=m.modelId===id); return delay(true);},
  setModelStatus:(id:number,status:'enabled'|'disabled')=>{const m=baseModels.find(x=>x.modelId===id); if(m)m.status=status; return delay(true);},
  getTokenConfig:()=>delay(tokenConfig),
  saveTokenConfig:(data:TokenConfig)=>{tokenConfig={...tokenConfig,...data,updatedAt:now()}; return delay(tokenConfig);},
  getTokenUsages:(params:PageRequest)=>delay(page(tokenUsages,params.page,params.pageSize)),
  resetUserTokens:(userId:number)=>{const u=tokenUsages.find(x=>x.userId===userId); if(u)u.todayTokens=0; return delay(true);},
  getAdminLogs:(params:PageRequest)=>delay(page(adminLogs,params.page,params.pageSize)),
  adminStats:()=>delay({ todayNewUserCount:128,todayPostCount:256,todayCommentCount:1234,todayLikeCount:3200,pendingReviewCount:18,pendingReportCount:32,dau:2680,retentionRate:42.6,conversionRate:18.8 }),
};


function realEndpoint(methodName: string, args: any[]): { method: 'get' | 'post' | 'put' | 'delete'; url: string; params?: any; data?: any } {
  const [a, b, c] = args;
  const maps: Record<string, () => { method: 'get' | 'post' | 'put' | 'delete'; url: string; params?: any; data?: any }> = {
    login: () => ({ method: 'post', url: '/auth/login', data: { account: a, password: b } }),
    register: () => ({ method: 'post', url: '/auth/register', data: a }),
    getUsers: () => ({ method: 'get', url: '/users', params: a }),
    updateUser: () => ({ method: 'put', url: '/users/me/profile', data: a }),
    getTags: () => ({ method: 'get', url: '/tags' }),
    getTopics: () => ({ method: 'get', url: '/topics', params: a }),
    getTopic: () => ({ method: 'get', url: `/topics/${a}` }),
    getCircles: () => ({ method: 'get', url: '/circles', params: a }),
    getCircle: () => ({ method: 'get', url: `/circles/${a}` }),
    createCircle: () => ({ method: 'post', url: '/circles', data: a }),
    joinCircle: () => ({ method: 'post', url: `/circles/${a}/join`, data: { reason: b } }),
    getPosts: () => ({ method: 'get', url: '/posts', params: a }),
    getPost: () => ({ method: 'get', url: `/posts/${a}` }),
    createPost: () => ({ method: 'post', url: '/posts', data: a }),
    updatePost: () => ({ method: 'put', url: `/posts/${a}`, data: b }),
    hidePost: () => ({ method: 'put', url: `/posts/${a}/hide` }),
    unhidePost: () => ({ method: 'put', url: `/posts/${a}/unhide` }),
    deletePost: () => ({ method: 'delete', url: `/posts/${a}` }),
    likePost: () => ({ method: 'post', url: `/posts/${a}/like` }),
    favoritePost: () => ({ method: 'post', url: `/posts/${a}/favorite` }),
    sharePost: () => ({ method: 'post', url: `/posts/${a}/share`, data: { channel: 'copy_link' } }),
    repostPost: () => ({ method: 'post', url: `/posts/${a}/repost`, data: { repostComment: b } }),
    getComments: () => ({ method: 'get', url: `/posts/${a}/comments` }),
    addComment: () => ({ method: 'post', url: '/comments', data: { postId: a, content: b } }),
    getMyComments: () => ({ method: 'get', url: '/users/me/comments', params: a }),
    deleteComment: () => ({ method: 'delete', url: `/comments/${a}` }),
    getCircleMembers: () => ({ method: 'get', url: `/circles/${a}/members`, params: b }),
    setMemberRole: () => ({ method: 'put', url: `/circles/${a}/members/${b}/role`, data: { role: c } }),
    muteMember: () => ({ method: 'put', url: `/circles/${a}/members/${b}/mute`, data: { duration: c, reason: args[3] } }),
    unmuteMember: () => ({ method: 'put', url: `/circles/${a}/members/${b}/unmute` }),
    removeMember: () => ({ method: 'delete', url: `/circles/${a}/members/${b}` }),
    getJoinRequests: () => ({ method: 'get', url: `/circles/${a}/join-requests` }),
    handleJoinRequest: () => ({ method: 'put', url: `/circles/join-requests/${a}/${b}` }),
    getAnnouncements: () => ({ method: 'get', url: `/circles/${a}/announcements` }),
    addAnnouncement: () => ({ method: 'post', url: `/circles/${a}/announcements`, data: { title: b, content: c } }),
    getChatMessages: () => ({ method: 'get', url: `/circles/${a}/chat/messages` }),
    sendChatMessage: () => ({ method: 'post', url: `/circles/${a}/chat/messages`, data: { content: b } }),
    getSummarySubscription: () => ({ method: 'get', url: `/circles/${a}/ai-summary/subscription` }),
    saveSummarySubscription: () => ({ method: 'put', url: `/circles/${a}/ai-summary/subscription`, data: b }),
    cancelSummarySubscription: () => ({ method: 'put', url: `/circles/${a}/ai-summary/subscription/cancel` }),
    getSummarySubscriptions: () => ({ method: 'get', url: '/users/me/circle-summary-subscriptions' }),
    getBanners: () => ({ method: 'get', url: '/banners', params: { position: 'home' } }),
    getHotRanks: () => ({ method: 'get', url: '/hot/ranks', params: { rankType: a, timeRange: b } }),
    getActivities: () => ({ method: 'get', url: '/activities', params: a }),
    getActivity: () => ({ method: 'get', url: `/activities/${a}` }),
    getOfficialAnnouncements: () => ({ method: 'get', url: '/announcements', params: a }),
    getNotificationList: () => ({ method: 'get', url: '/notifications' }),
    markNotificationRead: () => ({ method: 'put', url: `/notifications/${a}/read` }),
    getFollowingUsers: () => ({ method: 'get', url: '/following-center/users', params: a }),
    getFollowingCircles: () => ({ method: 'get', url: '/following-center/circles', params: a }),
    getFollowingTopics: () => ({ method: 'get', url: '/following-center/topics', params: a }),
    getFollowingTags: () => ({ method: 'get', url: '/following-center/tags', params: a }),
    getFollowingFeed: () => ({ method: 'get', url: '/following-center/feed', params: a }),
    search: () => ({ method: 'get', url: '/search', params: { keyword: a, type: b, ...(c || {}) } }),
    searchSuggest: () => ({ method: 'get', url: '/search/suggest', params: { keyword: a } }),
    getHotKeywords: () => ({ method: 'get', url: '/search/hot-keywords' }),
    getWorkspaceDashboard: () => ({ method: 'get', url: '/workspace/dashboard' }),
    getNotes: () => ({ method: 'get', url: '/workspace/notes', params: a }),
    getNote: () => ({ method: 'get', url: `/workspace/notes/${a}` }),
    saveNote: () => a?.noteId ? ({ method: 'put', url: `/workspace/notes/${a.noteId}`, data: a }) : ({ method: 'post', url: '/workspace/notes', data: a }),
    deleteNote: () => ({ method: 'delete', url: `/workspace/notes/${a}` }),
    publishNoteAsPost: () => ({ method: 'post', url: `/workspace/notes/${a}/publish-as-post` }),
    optimizeNote: () => ({ method: 'post', url: '/ai/notes/optimize', data: { content: a, goal: b } }),
    getKnowledgeBases: () => ({ method: 'get', url: '/workspace/knowledge-bases', params: a }),
    getKnowledgeBase: () => ({ method: 'get', url: `/workspace/knowledge-bases/${a}` }),
    saveKnowledgeBase: () => a?.knowledgeBaseId ? ({ method: 'put', url: `/workspace/knowledge-bases/${a.knowledgeBaseId}`, data: a }) : ({ method: 'post', url: '/workspace/knowledge-bases', data: a }),
    getKnowledgeBaseMembers: () => ({ method: 'get', url: `/workspace/knowledge-bases/${a}/members` }),
    inviteKbMember: () => ({ method: 'post', url: `/workspace/knowledge-bases/${a}/invite`, data: { userId: b, role: c } }),
    updateKbMemberRole: () => ({ method: 'put', url: `/workspace/knowledge-bases/${a}/members/${b}/role`, data: { role: c } }),
    removeKbMember: () => ({ method: 'delete', url: `/workspace/knowledge-bases/${a}/members/${b}` }),
    getModelConfig: () => ({ method: 'get', url: '/ai/model-config' }),
    saveModelConfig: () => ({ method: 'put', url: '/ai/model-config', data: a }),
    testModelConfig: () => ({ method: 'post', url: '/ai/model-config/test', data: a }),
    getAiSessions: () => ({ method: 'get', url: '/ai/chat/sessions' }),
    createAiSession: () => ({ method: 'post', url: '/ai/chat/sessions', data: { title: a } }),
    getAiMessages: () => ({ method: 'get', url: `/ai/chat/sessions/${a}` }),
    sendAiMessage: () => ({ method: 'post', url: `/ai/chat/sessions/${a}/messages`, data: { content: b } }),
    updateAiMessage: () => ({ method: 'put', url: `/ai/chat/messages/${a}`, data: { content: b } }),
    saveAiToNote: () => ({ method: 'post', url: `/ai/answers/${a}/save-to-note`, data: { title: b, knowledgeBaseId: c } }),
    topicAgent: () => ({ method: 'post', url: '/ai/topic-agent', data: { topicId: a, action: b, question: c } }),
    checkIn: () => ({ method: 'post', url: '/growth/check-in' }),
    getTasks: () => ({ method: 'get', url: '/growth/tasks', params: { type: a } }),
    claimTask: () => ({ method: 'post', url: `/growth/tasks/${a}/claim` }),
    getBadges: () => ({ method: 'get', url: `/growth/users/${a || 'me'}/badges` }),
    getRankings: () => ({ method: 'get', url: '/growth/rankings', params: { type: a, range: b } }),
    getBaseModels: () => ({ method: 'get', url: '/admin/ai/models', params: a }),
    saveBaseModel: () => a?.modelId ? ({ method: 'put', url: `/admin/ai/models/${a.modelId}`, data: a }) : ({ method: 'post', url: '/admin/ai/models', data: a }),
    setDefaultModel: () => ({ method: 'put', url: `/admin/ai/models/${a}/default` }),
    setModelStatus: () => ({ method: 'put', url: `/admin/ai/models/${a}/${b === 'enabled' ? 'enable' : 'disable'}` }),
    getTokenConfig: () => ({ method: 'get', url: '/admin/ai/token-config' }),
    saveTokenConfig: () => ({ method: 'put', url: '/admin/ai/token-config', data: a }),
    getTokenUsages: () => ({ method: 'get', url: '/admin/ai/token-usages', params: a }),
    resetUserTokens: () => ({ method: 'post', url: `/admin/ai/token-usages/${a}/reset-today` }),
    getAdminLogs: () => ({ method: 'get', url: '/admin/operation-logs', params: a }),
    adminStats: () => ({ method: 'get', url: '/admin/dashboard/stats' }),
  };
  return maps[methodName]?.() || { method: 'post', url: `/__rpc/${methodName}`, data: { args } };
}

function createRealApi<T extends Record<string, any>>(shape: T): T {
  return new Proxy(shape, {
    get(_target, prop) {
      const methodName = String(prop);
      if (methodName === 'delay') return mockApi.delay;
      return async (...args: any[]) => {
        const endpoint = realEndpoint(methodName, args);
        if (endpoint.method === 'get') return httpGet(endpoint.url, endpoint.params);
        if (endpoint.method === 'post') return httpPost(endpoint.url, endpoint.data);
        if (endpoint.method === 'put') return httpPut(endpoint.url, endpoint.data);
        return httpDelete(endpoint.url, endpoint.params);
      };
    },
  }) as T;
}

export const realApi = createRealApi(mockApi);
export const api = createApiAdapter(mockApi, realApi);
