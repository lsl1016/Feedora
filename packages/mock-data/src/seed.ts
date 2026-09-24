const svg = (text: string, bg: string) =>
  'data:image/svg+xml;utf8,' + encodeURIComponent(
    '<svg xmlns="http://www.w3.org/2000/svg" width="160" height="160"><rect width="100%" height="100%" rx="28" fill="' + bg + '"/><text x="50%" y="54%" dominant-baseline="middle" text-anchor="middle" fill="white" font-family="Arial" font-size="48" font-weight="700">' + text + '</text></svg>'
  );

const cover = (title: string, a: string, b: string) =>
  'data:image/svg+xml;utf8,' + encodeURIComponent(
    '<svg xmlns="http://www.w3.org/2000/svg" width="1200" height="480"><defs><linearGradient id="g" x1="0" x2="1"><stop stop-color="' + a + '"/><stop offset="1" stop-color="' + b + '"/></linearGradient></defs><rect width="100%" height="100%" fill="url(%23g)"/><circle cx="980" cy="110" r="150" fill="rgba(255,255,255,.12)"/><circle cx="1080" cy="350" r="210" fill="rgba(255,255,255,.08)"/><text x="70" y="250" fill="white" font-family="Arial" font-size="64" font-weight="700">' + title + '</text></svg>'
  );

export const seedUsers = [
  { userId: 1, account: 'linchen', nickname: 'Lin Chen', avatar: svg('LC','#2563eb'), bio: 'Go / AI / 数据平台工程师', level: 7, levelName: '高级创作者', role: 'admin', status: 'normal', points: 2680, experience: 7400, nextLevelExperience: 9000, badgeCount: 12, checkedInToday: false, continuousCheckInDays: 8, postCount: 36, commentCount: 128, followerCount: 1280, followingCount: 96, likeReceivedCount: 8200, createdAt: '2026-01-03 09:00:00', updatedAt: '2026-09-24 10:24:00' },
  { userId: 2, account: 'shirley', nickname: 'Shirley', avatar: svg('SH','#7c3aed'), bio: 'RAG / Agent / 知识工程', level: 6, levelName: '知识探索者', role: 'user', status: 'normal', points: 1980, experience: 6100, nextLevelExperience: 7000, badgeCount: 8, checkedInToday: true, continuousCheckInDays: 15, postCount: 28, commentCount: 86, followerCount: 940, followingCount: 120, likeReceivedCount: 6100, createdAt: '2026-02-11 09:00:00', updatedAt: '2026-09-24 09:30:00' },
  { userId: 3, account: 'haoyu', nickname: '皓宇', avatar: svg('HY','#0f766e'), bio: 'Frontend / Desktop / Wails', level: 5, levelName: '实践者', role: 'user', status: 'normal', points: 1420, experience: 4200, nextLevelExperience: 5200, badgeCount: 6, checkedInToday: false, continuousCheckInDays: 3, postCount: 19, commentCount: 64, followerCount: 620, followingCount: 80, likeReceivedCount: 3300, createdAt: '2026-03-08 09:00:00', updatedAt: '2026-09-23 21:10:00' },
];

export const seedTags = [
  { tagId:1, tagName:'Go', description:'Go 语言与后端工程', status:'enabled', useCount:2380, createdAt:'2026-01-01', updatedAt:'2026-09-24' },
  { tagId:2, tagName:'AI Agent', description:'Agent Runtime / Tool / MCP', status:'enabled', useCount:1860, createdAt:'2026-01-01', updatedAt:'2026-09-24' },
  { tagId:3, tagName:'RAG', description:'知识库与检索增强生成', status:'enabled', useCount:1210, createdAt:'2026-01-01', updatedAt:'2026-09-24' },
  { tagId:4, tagName:'桌面应用', description:'Wails / Electron / Tauri', status:'enabled', useCount:640, createdAt:'2026-01-01', updatedAt:'2026-09-24' },
  { tagId:5, tagName:'架构设计', description:'系统设计与工程实践', status:'enabled', useCount:3020, createdAt:'2026-01-01', updatedAt:'2026-09-24' },
  { tagId:6, tagName:'数据平台', description:'数据研发与平台工程', status:'enabled', useCount:980, createdAt:'2026-01-01', updatedAt:'2026-09-24' },
];

export const seedTopics = [
  { topicId:1, name:'AI 工程化', description:'从模型调用到 Runtime、Tool、Memory 与评测。', coverImage:cover('AI ENGINEERING','#2563eb','#7c3aed'), postCount:286, participantCount:4260, isOfficial:true, isRecommended:true, status:'enabled', createdAt:'2026-03-01', updatedAt:'2026-09-24' },
  { topicId:2, name:'Go 后端实践', description:'Go 服务端、高并发、可观测与工程化。', coverImage:cover('GO BACKEND','#0f766e','#0ea5e9'), postCount:412, participantCount:5960, isOfficial:false, isRecommended:true, status:'enabled', createdAt:'2026-03-01', updatedAt:'2026-09-24' },
  { topicId:3, name:'开发者桌面工具', description:'本地优先、桌面客户端与多端协同。', coverImage:cover('DESKTOP DEV','#111827','#2563eb'), postCount:138, participantCount:1880, isOfficial:false, isRecommended:true, status:'enabled', createdAt:'2026-05-01', updatedAt:'2026-09-24' },
];

export const seedCircles = [
  { circleId:1, name:'Agent Runtime 研究所', avatar:svg('AR','#2563eb'), description:'讨论 Agent Runtime、MCP、Skill、Memory 与协议设计。', category:'AI', tags:[{tagId:2,tagName:'AI Agent'},{tagId:5,tagName:'架构设计'}], ownerId:1, owner:seedUsers[0], joinType:'direct', postPermission:'all', memberCount:5260, postCount:482, featuredPostCount:38, isJoined:true, myRole:'owner', myStatus:'normal', isRecommended:true, status:'normal', rules:'技术讨论优先，禁止广告。', createdAt:'2026-04-01', updatedAt:'2026-09-24' },
  { circleId:2, name:'Go 工程实践', avatar:svg('GO','#0f766e'), description:'Go Web、并发、微服务、性能优化与代码设计。', category:'Backend', tags:[{tagId:1,tagName:'Go'}], ownerId:2, owner:seedUsers[1], joinType:'direct', postPermission:'all', memberCount:8420, postCount:960, featuredPostCount:72, isJoined:true, myRole:'member', myStatus:'normal', isRecommended:true, status:'normal', rules:'保持代码和结论可复现。', createdAt:'2026-02-01', updatedAt:'2026-09-24' },
  { circleId:3, name:'桌面生产力工具', avatar:svg('DT','#7c3aed'), description:'Wails、Tauri、Electron 与本地知识工作台。', category:'Desktop', tags:[{tagId:4,tagName:'桌面应用'}], ownerId:3, owner:seedUsers[2], joinType:'approval', postPermission:'all', memberCount:1880, postCount:226, featuredPostCount:19, isJoined:false, isRecommended:true, status:'normal', rules:'欢迎分享真实项目经验。', createdAt:'2026-06-01', updatedAt:'2026-09-24' },
];

const brief = (c:any) => ({circleId:c.circleId,name:c.name,avatar:c.avatar,memberCount:c.memberCount,postCount:c.postCount});
export const seedPosts = [
  { postId:101, postType:'original', authorId:1, author:seedUsers[0], title:'Feedora Desktop V1：从社区到本地 Knowledge IDE 的架构设计', content:'我们把 Feedora Desktop 定位为本地知识工作台、在线知识社区和 AI Agent 的组合。核心边界是：社区事实来自远端服务，本地 Workspace 由 SQLite 和文件系统管理，Agent Runtime 独立部署。', summary:'Web/Desktop 双宿主、Local-first Workspace、Wails Adapter 与 Agent Runtime 解耦。', images:[], tags:[{tagId:4,tagName:'桌面应用'},{tagId:5,tagName:'架构设计'}], topics:[{topicId:3,name:'开发者桌面工具'}], circleId:3, circle:brief(seedCircles[2]), visibility:'public', status:'published', isTop:true, isFeatured:true, isSelected:true, viewCount:12840, likeCount:684, commentCount:92, favoriteCount:438, shareCount:86, repostCount:32, hotScore:9860, liked:false, favorited:true, followedAuthor:false, createdAt:'2026-09-24 10:24:00', publishedAt:'2026-09-24 10:24:00', updatedAt:'2026-09-24 10:24:00' },
  { postId:102, postType:'original', authorId:2, author:seedUsers[1], title:'RAG 在个人知识库中的实践：索引、重排与引用回溯', content:'把 Markdown、PDF、图片和仓库文档统一变成可追溯资源，再用混合检索和重排构建个人知识工作流。', summary:'从文档解析到混合检索，再到引用回溯的一套可落地方案。', images:[], tags:[{tagId:3,tagName:'RAG'},{tagId:2,tagName:'AI Agent'}], topics:[{topicId:1,name:'AI 工程化'}], circleId:1, circle:brief(seedCircles[0]), visibility:'public', status:'published', isTop:false, isFeatured:true, isSelected:false, viewCount:9600, likeCount:512, commentCount:66, favoriteCount:390, shareCount:49, repostCount:18, hotScore:8420, liked:true, favorited:false, followedAuthor:true, createdAt:'2026-09-23 18:12:00', publishedAt:'2026-09-23 18:12:00', updatedAt:'2026-09-23 18:12:00' },
  { postId:103, postType:'original', authorId:3, author:seedUsers[2], title:'为什么我最终选择 Wails 做 Go 项目的桌面客户端', content:'当后端和本地能力已经是 Go 时，Wails 能让 React UI 与 Go 本地服务保持清晰边界。', summary:'Wails 在 Go 技术栈下的工程收益、限制与适用边界。', images:[], tags:[{tagId:1,tagName:'Go'},{tagId:4,tagName:'桌面应用'}], topics:[{topicId:3,name:'开发者桌面工具'}], circleId:3, circle:brief(seedCircles[2]), visibility:'public', status:'published', isTop:false, isFeatured:false, isSelected:false, viewCount:7340, likeCount:391, commentCount:54, favoriteCount:201, shareCount:36, repostCount:12, hotScore:7210, liked:false, favorited:false, followedAuthor:false, createdAt:'2026-09-22 14:50:00', publishedAt:'2026-09-22 14:50:00', updatedAt:'2026-09-22 14:50:00' },
  { postId:104, postType:'original', authorId:1, author:seedUsers[0], title:'Go 服务端工具系统为什么要拆 Registry / Scheduler / Executor', content:'工具发现、调度和执行三个职责如果混在一起，后续很难处理并发、权限、恢复和可观测。', summary:'三个层次分别解决工具发现、执行编排和实际运行。', images:[], tags:[{tagId:1,tagName:'Go'},{tagId:5,tagName:'架构设计'}], topics:[{topicId:2,name:'Go 后端实践'}], circleId:2, circle:brief(seedCircles[1]), visibility:'public', status:'published', isTop:false, isFeatured:false, isSelected:false, viewCount:6810, likeCount:344, commentCount:48, favoriteCount:188, shareCount:29, repostCount:9, hotScore:6530, liked:false, favorited:false, followedAuthor:false, createdAt:'2026-09-21 20:10:00', publishedAt:'2026-09-21 20:10:00', updatedAt:'2026-09-21 20:10:00' },
];

export const seedComments = [
  { commentId:1, postId:101, userId:2, user:seedUsers[1], content:'把 Desktop 当成 Web 的壳确实会限制后续本地能力，Port 这层很关键。', likeCount:24, liked:false, status:'normal', replies:[], createdAt:'2026-09-24 11:00:00', updatedAt:'2026-09-24 11:00:00' },
  { commentId:2, postId:101, userId:3, user:seedUsers[2], content:'期待 Repository + AI Context 这条链路。', likeCount:17, liked:true, status:'normal', replies:[], createdAt:'2026-09-24 11:18:00', updatedAt:'2026-09-24 11:18:00' },
];

export const seedKnowledgeBases = [
  { knowledgeBaseId:1, name:'Feedora Desktop', description:'桌面应用架构、功能设计与开发实现文档', ownerId:1, noteCount:12, memberCount:3, visibility:'collaborative', createdAt:'2026-09-01', updatedAt:'2026-09-24' },
  { knowledgeBaseId:2, name:'Go Backend', description:'Go 工程与高并发实践', ownerId:1, noteCount:18, memberCount:1, visibility:'private', createdAt:'2026-08-01', updatedAt:'2026-09-23' },
  { knowledgeBaseId:3, name:'RAG / AI', description:'RAG、Agent、MCP、Memory 学习沉淀', ownerId:1, noteCount:27, memberCount:2, visibility:'collaborative', createdAt:'2026-07-01', updatedAt:'2026-09-22' },
];

export const seedNotes = [
  { noteId:1, title:'Feedora Desktop V1 架构拆解', content:'# Feedora Desktop V1 架构拆解\n\n## 1. 整体架构\n\nFeedora Desktop 采用 **共享 Client Core + 双宿主 + 平台 Adapter**。\n\n## 2. 核心边界\n\n- Community：远端事实源\n- Workspace：Local-first\n- Repository：只读索引\n- Agent Runtime：独立部署', summary:'Web/Desktop 双宿主与核心模块边界。', knowledgeBaseId:1, knowledgeBaseName:'Feedora Desktop', ownerId:1, status:'normal', createdAt:'2026-09-20', updatedAt:'2026-09-24 10:24:00' },
  { noteId:2, title:'Wails Bridge 设计要点', content:'# Wails Bridge\n\nReact Feature 不能直接依赖 Wails Binding，必须通过 Port 和 Adapter。', summary:'Wails Binding 与业务层的隔离方式。', knowledgeBaseId:1, knowledgeBaseName:'Feedora Desktop', ownerId:1, status:'normal', createdAt:'2026-09-21', updatedAt:'2026-09-24 09:42:00' },
  { noteId:3, title:'RAG 检索链路梳理', content:'# RAG 检索链路\n\n解析 → Chunk → FTS/Vector → Merge → Rerank → Citation。', summary:'个人知识库的混合检索与重排流程。', knowledgeBaseId:3, knowledgeBaseName:'RAG / AI', ownerId:1, status:'normal', createdAt:'2026-09-18', updatedAt:'2026-09-23 20:00:00' },
  { noteId:4, title:'Repository Indexing 设计', content:'# Repository Indexing\n\n首版只读扫描 Git 仓库，记录 file / symbol / commit / diff。', summary:'Git 仓库只读索引设计。', knowledgeBaseId:1, knowledgeBaseName:'Feedora Desktop', ownerId:1, status:'normal', createdAt:'2026-09-19', updatedAt:'2026-09-22 18:00:00' },
];

export const seedNotifications = [
  { notificationId:1, title:'Shirley 评论了你的帖子', content:'Port 这层很关键，尤其是 Web/Desktop 的能力差异。', category:'interaction', targetUrl:'/posts/101', readStatus:'unread', createdAt:'2026-09-24 11:00:00' },
  { notificationId:2, title:'你的帖子被收藏', content:'“Feedora Desktop V1 架构设计”新增 12 次收藏。', category:'interaction', targetUrl:'/posts/101', readStatus:'unread', createdAt:'2026-09-24 10:48:00' },
  { notificationId:3, title:'圈子有新的精华内容', content:'Agent Runtime 研究所新增一篇 Runtime Command Queue 复盘。', category:'circle', targetUrl:'/circles/1', readStatus:'unread', createdAt:'2026-09-24 09:30:00' },
  { notificationId:4, title:'索引任务完成', content:'feedora 仓库索引完成，共处理 128 个文件。', category:'system', targetUrl:'/workspace', readStatus:'read', createdAt:'2026-09-24 09:12:00' },
];

export const seedTasks = [
  { taskId:1, title:'完善个人资料', description:'补充头像、简介和技术方向。', type:'newbie', rewardPoints:50, targetValue:1, currentValue:1, status:'done', actionText:'查看资料', actionUrl:'/users/me' },
  { taskId:2, title:'发布一篇技术内容', description:'今天发布一篇原创技术帖子。', type:'daily', rewardPoints:30, targetValue:1, currentValue:0, status:'todo', actionText:'去发布', actionUrl:'/posts/create' },
  { taskId:3, title:'沉淀 5 篇知识笔记', description:'在工作空间中创建并维护 5 篇笔记。', type:'growth', rewardPoints:120, targetValue:5, currentValue:4, status:'todo', actionText:'去写笔记', actionUrl:'/workspace/notes' },
];

export const seedActivities = [
  { activityId:1, type:'topic', title:'桌面生产力工具周', description:'分享你正在做的桌面开发工具、架构和真实踩坑。', coverImage:cover('DESKTOP WEEK','#111827','#2563eb'), topicId:3, topicName:'开发者桌面工具', startAt:'2026-09-20', endAt:'2026-10-05', participantCount:842, submissionCount:126, status:'ongoing' },
  { activityId:2, type:'vote', title:'2026 开发者最想要的 AI 能力', description:'选择你最希望知识社区优先建设的 AI 能力。', coverImage:cover('AI VOTE','#7c3aed','#2563eb'), startAt:'2026-09-22', endAt:'2026-09-30', participantCount:1520, submissionCount:1520, status:'ongoing', voteConfig:{question:'你最希望优先增强哪项能力？',multiple:false,voted:false,selectedOptionIds:[],options:[{optionId:1,text:'本地知识库',voteCount:620,percent:41},{optionId:2,text:'代码仓库问答',voteCount:510,percent:34},{optionId:3,text:'自动写作',voteCount:390,percent:25}]} },
  { activityId:3, type:'checkin', title:'连续 21 天技术笔记挑战', description:'每天沉淀一条可复用的技术知识。', coverImage:cover('21 DAYS','#0f766e','#0ea5e9'), startAt:'2026-09-10', endAt:'2026-10-01', participantCount:680, submissionCount:4300, status:'ongoing', checkInConfig:{totalDays:21,currentDays:8,todayChecked:false,todayTask:'写下今天最重要的一条技术收获'} },
];

export const seedAnnouncements = [
  { announcementId:1, title:'Feedora Desktop V1 开始内测', summary:'本地知识工作台与社区联动能力进入第一轮内测。', content:'Feedora Desktop V1 已进入内测阶段。首批能力包括社区阅读、Workspace、Repository、统一搜索与 AI Panel。', publisherName:'Feedora Team', status:'published', publishedAt:'2026-09-24 08:30:00', createdAt:'2026-09-24', updatedAt:'2026-09-24' },
  { announcementId:2, title:'社区内容规范更新', summary:'新增 AI 生成内容标识与敏感代码处理建议。', content:'发布由 AI 辅助生成的内容时，请确保不包含公司私有代码、密钥和未授权资料。', publisherName:'Feedora Team', status:'published', publishedAt:'2026-09-20 12:00:00', createdAt:'2026-09-20', updatedAt:'2026-09-20' },
];

export const seedBanners = [
  { bannerId:1, title:'Feedora Desktop V1', imageUrl:cover('FEEDORA DESKTOP','#2563eb','#7c3aed'), targetUrl:'/workspace', position:'home', weight:100, status:'enabled', createdAt:'2026-09-24' },
];

export const seedAiSessions = [{ sessionId:1, title:'Feedora Desktop 架构总结', userId:1, createdAt:'2026-09-24 09:00:00', updatedAt:'2026-09-24 10:24:00' }];
export const seedAiMessages = [
  { messageId:1, sessionId:1, role:'user', content:'帮我总结 Feedora Desktop 的架构边界。', editable:false, createdAt:'2026-09-24 10:20:00' },
  { messageId:2, sessionId:1, role:'assistant', content:'核心可以归纳为三点：共享 Client Core；Web/Desktop 双宿主；本地能力通过 Wails Adapter 接入，Agent Runtime 保持独立。', editable:true, createdAt:'2026-09-24 10:20:05' },
];

export const seedModelConfig = { configId:1, userId:1, provider:'deepseek', baseUrl:'https://api.deepseek.com', apiKeyMasked:'sk-************demo', modelName:'deepseek-chat', temperature:0.7, maxTokens:4096, status:'enabled', createdAt:'2026-09-01', updatedAt:'2026-09-24' };

export const seedBaseModels = [
  { modelId:1, provider:'deepseek', providerName:'DeepSeek', modelName:'deepseek-chat', baseUrl:'https://api.deepseek.com', isDefault:true, status:'enabled', createdAt:'2026-09-01', updatedAt:'2026-09-24' },
  { modelId:2, provider:'qwen', providerName:'Qwen', modelName:'qwen-plus', baseUrl:'https://dashscope.aliyuncs.com/compatible-mode/v1', isDefault:false, status:'enabled', createdAt:'2026-09-01', updatedAt:'2026-09-24' },
];

export const seedRepositories = [
  { id:'repo-feedora', name:'feedora', path:'~/workspace/feedora', branch:'main', status:'ready', files:128, lastIndexedAt:'2026-09-24 09:42' },
  { id:'repo-runtime', name:'react-base-service', path:'~/workspace/react-base-service', branch:'main', status:'ready', files:264, lastIndexedAt:'2026-09-23 21:30' },
];

export const seedRepositoryFiles = [
  { repositoryId:'repo-feedora', path:'apps/desktop/main.go', language:'go', content:'package main\n\n// Feedora Desktop bootstrap.\nfunc main() {\n    // Wails application starts here.\n}\n' },
  { repositoryId:'repo-feedora', path:'packages/app-core/src/ports/search.ts', language:'typescript', content:'export interface SearchPort {\n  search(query: SearchQuery): Promise<SearchResult[]>;\n}\n' },
  { repositoryId:'repo-feedora', path:'docs/client-platform-architecture.md', language:'markdown', content:'# Feedora Client Platform\n\nWeb 与 Desktop 是两个宿主，共享 Client Core 和 Feature。\n' },
];

export const seedSearchResults = [
  { id:'note-2', source:'note', title:'Wails Bridge 设计要点', summary:'React Feature 不直接依赖 Wails Binding。', location:'Feedora Desktop / 技术笔记', score:0.98 },
  { id:'code-1', source:'code', title:'SearchPort', summary:'统一搜索的 Port 接口定义。', location:'feedora/packages/app-core/src/ports/search.ts:1', score:0.95 },
  { id:'resource-1', source:'resource', title:'Feedora Desktop V1 产品与架构设计', summary:'产品定位、Wails 集成、数据模型与 UI 冻结方案。', location:'Workspace / resources / product_prd.pdf', score:0.91 },
  { id:'community-101', source:'community', title:'Feedora Desktop V1：从社区到本地 Knowledge IDE 的架构设计', summary:'Web/Desktop 双宿主、Local-first Workspace。', location:'Community / Desktop', score:0.88 },
];
