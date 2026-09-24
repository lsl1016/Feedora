import {
  seedActivities,
  seedAiMessages,
  seedAiSessions,
  seedAnnouncements,
  seedBaseModels,
  seedBanners,
  seedCircles,
  seedComments,
  seedKnowledgeBases,
  seedModelConfig,
  seedNotes,
  seedNotifications,
  seedPosts,
  seedTags,
  seedTasks,
  seedTopics,
  seedUsers,
} from './seed';

type Fn = (...args: any[]) => Promise<any>;
const clone = <T,>(value: T): T => JSON.parse(JSON.stringify(value));
const wait = <T,>(value: T, ms = 40) => new Promise<T>((resolve) => setTimeout(() => resolve(clone(value)), ms));
const page = (list: any[], q: any = {}) => {
  const p = Number(q?.page || 1), size = Number(q?.pageSize || 20);
  return { list: list.slice((p - 1) * size, p * size), total: list.length, page: p, pageSize: size };
};

export function createMockFeedoraApi(): Record<string, Fn> {
  const state:any = {
    users: clone(seedUsers), tags: clone(seedTags), topics: clone(seedTopics), circles: clone(seedCircles),
    posts: clone(seedPosts), comments: clone(seedComments), notes: clone(seedNotes), kbs: clone(seedKnowledgeBases),
    notifications: clone(seedNotifications), tasks: clone(seedTasks), activities: clone(seedActivities),
    announcements: clone(seedAnnouncements), aiSessions: clone(seedAiSessions), aiMessages: clone(seedAiMessages),
    modelConfig: clone(seedModelConfig), baseModels: clone(seedBaseModels),
    joinRequests: [{requestId:1,circleId:3,userId:2,user:clone(seedUsers[1]),reason:'希望参与桌面端工程讨论',status:'pending',createdAt:'2026-09-24 08:50:00'}],
    circleAnnouncements: [{announcementId:1,circleId:1,title:'本周讨论：Runtime Command Queue',content:'欢迎提交你的设计和实现经验。',publisher:clone(seedUsers[0]),createdAt:'2026-09-24 08:00:00'}],
    chat: [{messageId:1,circleId:1,senderId:2,sender:clone(seedUsers[1]),content:'今天有人继续研究 Protocol V4 吗？',messageType:'text',status:'normal',createdAt:'2026-09-24 09:30:00'}],
    summarySubscriptions: {} as Record<number,any>,
    kbMembers: {
      1: [
        {id:1,knowledgeBaseId:1,userId:1,user:clone(seedUsers[0]),role:'owner',status:'active',joinedAt:'2026-09-01'},
        {id:2,knowledgeBaseId:1,userId:2,user:clone(seedUsers[1]),role:'editor',status:'active',joinedAt:'2026-09-05'},
        {id:3,knowledgeBaseId:1,userId:3,user:clone(seedUsers[2]),role:'viewer',status:'active',joinedAt:'2026-09-06'},
      ],
      2: [{id:4,knowledgeBaseId:2,userId:1,user:clone(seedUsers[0]),role:'owner',status:'active',joinedAt:'2026-08-01'}],
      3: [{id:5,knowledgeBaseId:3,userId:1,user:clone(seedUsers[0]),role:'owner',status:'active',joinedAt:'2026-07-01'}],
    },
    tokenConfig: { configId:1,dailyMaxTokens:2000000,userDailyMaxTokens:100000,singleConversationMaxTokens:20000,overLimitStrategy:'warn',updatedAt:'2026-09-24' },
  };
  let nextId = 3000;
  const currentUser = () => state.users[0];
  const userSummary = (u:any) => ({userId:u.userId,nickname:u.nickname,avatar:u.avatar,bio:u.bio,level:u.level,levelName:u.levelName});

  const methods:Record<string,Fn> = {
    login: async () => wait({token:'mock-feedora-token',user:currentUser()}),
    register: async (input:any) => wait({userId:++nextId, ...input}),
    getUsers: async (q:any) => wait(page(state.users,q)),
    getUser: async (id:number) => wait(state.users.find((x:any)=>x.userId===Number(id)) || currentUser()),
    updateUser: async (input:any) => { Object.assign(currentUser(),input); return wait(currentUser()); },

    getTags: async () => wait(state.tags),
    getTag: async (id:number) => wait(state.tags.find((x:any)=>x.tagId===Number(id)) || state.tags[0]),
    getTopics: async (q:any={}) => wait(page(state.topics.filter((x:any)=>q.type!=='official'||x.isOfficial),q)),
    getTopic: async (id:number) => wait(state.topics.find((x:any)=>x.topicId===Number(id)) || state.topics[0]),

    getCircles: async (q:any={}) => {
      let rows=[...state.circles];
      if(q.scope==='recommended') rows=rows.filter((x:any)=>x.isRecommended);
      if(q.scope==='joined') rows=rows.filter((x:any)=>x.isJoined);
      if(q.scope==='created') rows=rows.filter((x:any)=>x.ownerId===currentUser().userId);
      if(q.keyword) rows=rows.filter((x:any)=>x.name.includes(q.keyword)||x.description.includes(q.keyword));
      return wait(page(rows,q));
    },
    getCircle: async (id:number) => wait(state.circles.find((x:any)=>x.circleId===Number(id)) || state.circles[0]),
    createCircle: async (input:any) => {
      const c={circleId:++nextId,name:input.name||'新圈子',avatar:currentUser().avatar,description:input.description||'',category:input.category||'General',tags:state.tags.filter((t:any)=>(input.tagIds||[]).includes(t.tagId)),ownerId:1,owner:userSummary(currentUser()),joinType:input.joinType||'direct',postPermission:'all',memberCount:1,postCount:0,featuredPostCount:0,isJoined:true,myRole:'owner',myStatus:'normal',isRecommended:false,status:'normal',rules:input.rules||'',createdAt:'2026-09-24',updatedAt:'2026-09-24'}; state.circles.unshift(c); return wait(c);
    },
    joinCircle: async (id:number) => { const c=state.circles.find((x:any)=>x.circleId===Number(id)); if(c){ if(c.joinType==='approval') return wait({pending:true}); c.isJoined=true;c.memberCount++; } return wait({pending:false}); },
    leaveCircle: async (id:number) => { const c=state.circles.find((x:any)=>x.circleId===Number(id)); if(c)c.isJoined=false; return wait({success:true}); },

    getPosts: async (q:any={}) => {
      let rows=state.posts.filter((x:any)=>x.status!=='deleted');
      if(q.status&&q.status!=='all') rows=rows.filter((x:any)=>x.status===q.status);
      if(q.tagId) rows=rows.filter((x:any)=>x.tags.some((t:any)=>t.tagId===Number(q.tagId)));
      if(q.topicId) rows=rows.filter((x:any)=>x.topics.some((t:any)=>t.topicId===Number(q.topicId)));
      if(q.circleId) rows=rows.filter((x:any)=>x.circleId===Number(q.circleId));
      if(q.authorId) rows=rows.filter((x:any)=>x.authorId===Number(q.authorId));
      if(q.keyword) rows=rows.filter((x:any)=>x.title.includes(q.keyword)||x.summary.includes(q.keyword));
      if(q.sort==='hot'||q.feedType==='hot') rows.sort((a:any,b:any)=>b.hotScore-a.hotScore);
      else rows.sort((a:any,b:any)=>String(b.createdAt).localeCompare(String(a.createdAt)));
      return wait(page(rows,q));
    },
    getPost: async (id:number) => wait(state.posts.find((x:any)=>x.postId===Number(id)) || state.posts[0]),
    createPost: async (input:any) => {
      const p={postId:++nextId,postType:'original',authorId:1,author:userSummary(currentUser()),title:input.title,content:input.content,summary:String(input.content||'').slice(0,100),images:input.images||[],tags:state.tags.filter((t:any)=>(input.tagIds||[]).includes(t.tagId)),topics:state.topics.filter((t:any)=>(input.topicIds||[]).includes(t.topicId)).map((t:any)=>({topicId:t.topicId,name:t.name})),circleId:input.circleId,visibility:input.visibility||'public',status:input.publishMode==='draft'?'draft':'published',isTop:false,isFeatured:false,isSelected:false,viewCount:0,likeCount:0,commentCount:0,favoriteCount:0,shareCount:0,repostCount:0,hotScore:0,liked:false,favorited:false,followedAuthor:false,createdAt:'2026-09-24 12:00:00',publishedAt:'2026-09-24 12:00:00',updatedAt:'2026-09-24 12:00:00'}; state.posts.unshift(p); return wait(p);
    },
    updatePost: async (id:number,input:any) => { const p=state.posts.find((x:any)=>x.postId===Number(id)); if(p)Object.assign(p,input,{updatedAt:'2026-09-24 12:10:00'}); return wait(p); },
    hidePost: async (id:number) => { const p=state.posts.find((x:any)=>x.postId===Number(id)); if(p)p.status='hidden'; return wait(p); },
    unhidePost: async (id:number) => { const p=state.posts.find((x:any)=>x.postId===Number(id)); if(p)p.status='published'; return wait(p); },
    deletePost: async (id:number) => { const p=state.posts.find((x:any)=>x.postId===Number(id)); if(p)p.status='deleted'; return wait({success:true}); },
    likePost: async (id:number) => { const p=state.posts.find((x:any)=>x.postId===Number(id)); if(p){p.liked=!p.liked;p.likeCount+=p.liked?1:-1;} return wait(p); },
    favoritePost: async (id:number) => { const p=state.posts.find((x:any)=>x.postId===Number(id)); if(p){p.favorited=!p.favorited;p.favoriteCount+=p.favorited?1:-1;} return wait(p); },
    sharePost: async (id:number) => { const p=state.posts.find((x:any)=>x.postId===Number(id)); if(p)p.shareCount++; return wait({success:true}); },
    repostPost: async (id:number) => { const p=state.posts.find((x:any)=>x.postId===Number(id)); if(p)p.repostCount++; return wait({success:true}); },

    getComments: async (postId:number,q:any={}) => wait(page(state.comments.filter((x:any)=>x.postId===Number(postId)&&x.status!=='deleted'),q)),
    addComment: async (postId:number,content:string) => { const c={commentId:++nextId,postId:Number(postId),userId:1,user:userSummary(currentUser()),content,likeCount:0,liked:false,status:'normal',replies:[],createdAt:'2026-09-24 12:20:00',updatedAt:'2026-09-24 12:20:00'}; state.comments.unshift(c); const p=state.posts.find((x:any)=>x.postId===Number(postId)); if(p)p.commentCount++; return wait(c); },
    getMyComments: async (q:any={}) => wait(page(state.comments.map((c:any)=>({...c,postTitle:state.posts.find((p:any)=>p.postId===c.postId)?.title||''})),q)),
    deleteComment: async (id:number) => { const c=state.comments.find((x:any)=>x.commentId===Number(id)); if(c)c.status='deleted'; return wait({success:true}); },

    getCircleMembers: async (circleId:number,q:any={}) => {
      const c=state.circles.find((x:any)=>x.circleId===Number(circleId));
      const rows=state.users.map((u:any,i:number)=>({id:circleId*100+i,circleId,userId:u.userId,user:userSummary(u),role:u.userId===c?.ownerId?'owner':'member',status:'normal',joinedAt:'2026-09-01',updatedAt:'2026-09-24'}));
      return wait(page(rows,q));
    },
    setMemberRole: async () => wait({success:true}),
    muteMember: async () => wait({success:true}),
    unmuteMember: async () => wait({success:true}),
    removeMember: async () => wait({success:true}),
    getJoinRequests: async (circleId:number) => wait(state.joinRequests.filter((x:any)=>x.circleId===Number(circleId))),
    handleJoinRequest: async (_circleId:number,requestId:number,status:string) => { const r=state.joinRequests.find((x:any)=>x.requestId===Number(requestId)); if(r)r.status=status; return wait(r); },
    getAnnouncements: async (circleId:number) => wait(state.circleAnnouncements.filter((x:any)=>x.circleId===Number(circleId))),
    addAnnouncement: async (circleId:number,input:any) => { const a={announcementId:++nextId,circleId:Number(circleId),title:input.title,content:input.content,publisher:userSummary(currentUser()),createdAt:'2026-09-24 12:00:00'};state.circleAnnouncements.unshift(a);return wait(a);},
    getChatMessages: async (circleId:number) => wait(state.chat.filter((x:any)=>x.circleId===Number(circleId))),
    sendChatMessage: async (circleId:number,content:string) => { const m={messageId:++nextId,circleId:Number(circleId),senderId:1,sender:userSummary(currentUser()),content,messageType:'text',status:'normal',createdAt:'2026-09-24 12:00:00'};state.chat.push(m);return wait(m);},
    getSummarySubscription: async (circleId:number) => wait(state.summarySubscriptions[circleId]||{subscriptionId:circleId,userId:1,circleId:Number(circleId),enabled:true,frequency:'daily',pushTime:'09:00',channels:['notification'],contentScopes:['post','chat'],createdAt:'2026-09-24',updatedAt:'2026-09-24'}),
    saveSummarySubscription: async (circleId:number,input:any) => { const s={subscriptionId:Number(circleId),userId:1,circleId:Number(circleId),...input,createdAt:'2026-09-24',updatedAt:'2026-09-24'};state.summarySubscriptions[circleId]=s;return wait(s);},
    cancelSummarySubscription: async (circleId:number) => { delete state.summarySubscriptions[circleId]; return wait({success:true}); },

    getHotRanks: async (type:string) => {
      if(type==='circle') return wait(state.circles.map((c:any,i:number)=>({rank:i+1,circleId:c.circleId,name:c.name,memberCount:c.memberCount,postCount:c.postCount,featuredPostCount:c.featuredPostCount,hotScore:9000-i*700})));
      if(type==='topic') return wait(state.topics.map((t:any,i:number)=>({rank:i+1,topicId:t.topicId,name:t.name,participantCount:t.participantCount,postCount:t.postCount,hotScore:8800-i*620})));
      return wait(state.posts.slice(0,10).map((p:any,i:number)=>({rank:i+1,postId:p.postId,title:p.title,authorName:p.author.nickname,hotScore:p.hotScore,likeCount:p.likeCount,commentCount:p.commentCount})));
    },

    getNotificationList: async () => wait(state.notifications),
    markNotificationRead: async (id:number) => { const n=state.notifications.find((x:any)=>x.notificationId===Number(id));if(n)n.readStatus='read';return wait(n);},
    getUnreadCount: async () => wait(state.notifications.filter((x:any)=>x.readStatus==='unread').length),

    search: async (keyword:string,type:string='all',q:any={}) => {
      const k=(keyword||'').toLowerCase();
      let rows:any[]=[];
      if(type==='all'||type==='post') rows.push(...state.posts.filter((x:any)=>!k||x.title.toLowerCase().includes(k)||x.summary.toLowerCase().includes(k)));
      if(type==='all'||type==='user') rows.push(...state.users.filter((x:any)=>!k||x.nickname.toLowerCase().includes(k)||x.bio.toLowerCase().includes(k)));
      if(type==='all'||type==='topic') rows.push(...state.topics.filter((x:any)=>!k||x.name.toLowerCase().includes(k)));
      if(type==='all'||type==='circle') rows.push(...state.circles.filter((x:any)=>!k||x.name.toLowerCase().includes(k)));
      return wait(page(rows,q));
    },
    searchSuggest: async (keyword:string) => wait(['Wails Bridge','Agent Runtime','RAG 检索'].filter(x=>x.toLowerCase().includes((keyword||'').toLowerCase()))),
    getHotKeywords: async () => wait(['Agent Runtime','Wails','RAG','Go 并发','MCP']),

    checkIn: async () => { currentUser().checkedInToday=true;currentUser().points+=20;return wait({points:20}); },
    getTasks: async (type:string) => wait(state.tasks.filter((x:any)=>x.type===type)),
    claimTask: async (id:number) => { const t=state.tasks.find((x:any)=>x.taskId===Number(id));if(t)t.status='claimed';return wait(t); },
    getRankings: async (type:string) => {
      if(type==='circle') return wait(state.circles.map((c:any,i:number)=>({rank:i+1,targetId:c.circleId,targetType:'circle',name:c.name,avatar:c.avatar,memberCount:c.memberCount,postCount:c.postCount,score:9800-i*600})));
      return wait(state.users.map((u:any,i:number)=>({rank:i+1,targetId:u.userId,targetType:'user',name:u.nickname,avatar:u.avatar,level:u.level,levelName:u.levelName,points:u.points,postCount:u.postCount,likeReceivedCount:u.likeReceivedCount,score:9900-i*700,isCurrentUser:u.userId===1})));
    },

    getMyPosts: async (q:any={}) => wait(page(state.posts.filter((x:any)=>x.authorId===1),q)),
    getMyLikedPosts: async (q:any={}) => wait(page(state.posts.filter((x:any)=>x.liked),q)),
    getMyFavoritePosts: async (q:any={}) => wait(page(state.posts.filter((x:any)=>x.favorited),q)),
    followUser: async () => wait({followed:true}),
    unfollowUser: async () => wait({followed:false}),
    getUserFollowState: async () => wait({followed:false}),
    getFollowingUsers: async (q:any={}) => wait(page(state.users.slice(1),q)),
    getFollowers: async (_id:number,q:any={}) => wait(page(state.users,q)),
    getFollowingFeed: async (q:any={}) => wait(page(state.posts.map((p:any)=>({feedId:p.postId,feedType:'post',title:p.title,summary:p.summary,targetId:p.postId,targetUrl:'/posts/'+p.postId,sourceName:p.author.nickname,sourceAvatar:p.author.avatar,createdAt:p.createdAt})),q)),
    getFollowingCircles: async (q:any={}) => wait(page(state.circles.filter((x:any)=>x.isJoined),q)),
    getFollowingTopics: async (q:any={}) => wait(page(state.topics,q)),
    getFollowingTags: async (q:any={}) => wait(page(state.tags,q)),

    getActivities: async (q:any={}) => wait(page(state.activities.filter((x:any)=>q.type==='all'||!q.type||x.type===q.type),q)),
    getActivity: async (id:number) => wait(state.activities.find((x:any)=>x.activityId===Number(id)) || state.activities[0]),
    getOfficialAnnouncements: async (q:any={}) => wait(page(state.announcements,q)),
    getBanners: async () => wait(seedBanners),

    getWorkspaceDashboard: async () => wait({noteCount:state.notes.length,knowledgeBaseCount:state.kbs.length,aiChatCount:state.aiSessions.length,collaborationMemberCount:3}),
    getNotes: async (q:any={}) => { let rows=state.notes.filter((x:any)=>x.status!=='deleted'); const kbId=q.kbId??q.knowledgeBaseId;if(kbId)rows=rows.filter((x:any)=>x.knowledgeBaseId===Number(kbId));if(q.keyword)rows=rows.filter((x:any)=>x.title.includes(q.keyword));return wait(page(rows,q));},
    getNote: async (id:number) => wait(state.notes.find((x:any)=>x.noteId===Number(id)) || state.notes[0]),
    saveNote: async (input:any) => {
      let n=state.notes.find((x:any)=>x.noteId===Number(input.noteId));
      const kb=state.kbs.find((x:any)=>x.knowledgeBaseId===Number(input.knowledgeBaseId));
      if(n) Object.assign(n,input,{summary:String(input.content||'').slice(0,80),knowledgeBaseName:kb?.name,updatedAt:'2026-09-24 12:00:00'});
      else { n={noteId:++nextId,title:input.title||'未命名笔记',content:input.content||'',summary:String(input.content||'').slice(0,80),knowledgeBaseId:input.knowledgeBaseId,knowledgeBaseName:kb?.name,ownerId:1,status:'normal',createdAt:'2026-09-24 12:00:00',updatedAt:'2026-09-24 12:00:00'};state.notes.unshift(n);}
      return wait(n);
    },
    deleteNote: async (id:number) => { const n=state.notes.find((x:any)=>x.noteId===Number(id));if(n)n.status='deleted';return wait({success:true});},
    optimizeNote: async (content:string,goal:string) => wait('## AI 优化建议\n\n- 先说明背景和目标\n- 将核心模块按职责拆分\n- 增加数据流和失败路径\n- 最后给出演进计划\n\n原文要点：'+String(content||'').slice(0,120)+'\n\n优化目标：'+goal),
    getKnowledgeBases: async (q:any={}) => wait(page(state.kbs,q)),
    getKnowledgeBase: async (id:number) => wait(state.kbs.find((x:any)=>x.knowledgeBaseId===Number(id)) || state.kbs[0]),
    saveKnowledgeBase: async (input:any) => { const kb={knowledgeBaseId:++nextId,name:input.name,description:input.description||'',ownerId:1,noteCount:0,memberCount:1,visibility:input.visibility||'private',createdAt:'2026-09-24',updatedAt:'2026-09-24'};state.kbs.unshift(kb);return wait(kb);},
    getKnowledgeBaseMembers: async (id:number) => wait(state.kbMembers[id]||[]),
    inviteKbMember: async (id:number,userId:number,role:string) => { const user=state.users.find((x:any)=>x.userId===Number(userId)); const m={id:++nextId,knowledgeBaseId:Number(id),userId:Number(userId),user:userSummary(user),role,status:'active',joinedAt:'2026-09-24'};(state.kbMembers[id]||(state.kbMembers[id]=[])).push(m);return wait(m);},
    updateKbMemberRole: async (id:number,userId:number,role:string) => { const m=(state.kbMembers[id]||[]).find((x:any)=>x.userId===Number(userId));if(m)m.role=role;return wait(m);},
    removeKbMember: async (id:number,userId:number) => { state.kbMembers[id]=(state.kbMembers[id]||[]).filter((x:any)=>x.userId!==Number(userId));return wait({success:true});},

    getAiSessions: async () => wait(state.aiSessions),
    createAiSession: async (title:string) => { const s={sessionId:++nextId,title:title||'新会话',userId:1,createdAt:'2026-09-24 12:00:00',updatedAt:'2026-09-24 12:00:00'};state.aiSessions.unshift(s);return wait(s);},
    getAiMessages: async (sessionId:number) => wait(state.aiMessages.filter((x:any)=>x.sessionId===Number(sessionId))),
    sendAiMessage: async (sessionId:number,content:string) => { const u={messageId:++nextId,sessionId:Number(sessionId),role:'user',content,editable:false,createdAt:'2026-09-24 12:00:00'};const a={messageId:++nextId,sessionId:Number(sessionId),role:'assistant',content:'基于 Mock 数据，我建议把这个问题拆成：产品边界、Client Core、平台 Adapter、数据流和演进计划五部分。你还可以把社区帖子、笔记和仓库文件同时加入 Context。',editable:true,createdAt:'2026-09-24 12:00:02'};state.aiMessages.push(u,a);return wait([u,a]);},
    updateAiMessage: async (id:number,content:string) => { const m=state.aiMessages.find((x:any)=>x.messageId===Number(id));if(m)m.content=content;return wait(m);},
    saveAiToNote: async (messageId:number,title:string,kbId?:number) => { const m=state.aiMessages.find((x:any)=>x.messageId===Number(messageId));return methods.saveNote({title,content:m?.content||'',knowledgeBaseId:kbId});},
    topicAgent: async (_topicId:number,action:string,question?:string) => wait({answer:'这是 Mock Agent 的回答。当前动作：'+action+(question?'，问题：'+question:'')+'。建议先阅读高质量帖子，再把结论沉淀进 Workspace。'}),
    getModelConfig: async () => wait(state.modelConfig),
    saveModelConfig: async (input:any) => { Object.assign(state.modelConfig,input,{apiKeyMasked:input.apiKey?'sk-************mock':state.modelConfig.apiKeyMasked,updatedAt:'2026-09-24'});return wait(state.modelConfig);},
    testModelConfig: async () => wait({ok:true}),

    adminGetUsers: async (q:any={}) => wait(page(state.users,q)),
    adminGetPosts: async (q:any={}) => wait(page(state.posts.filter((x:any)=>!q.status||x.status===q.status),q)),
    adminGetComments: async (q:any={}) => wait(page(state.comments,q)),
    adminGetTags: async () => wait(state.tags),
    adminCreateTag: async (input:any) => { const t={tagId:++nextId,tagName:input.name,description:input.description||'',status:'enabled',useCount:0,createdAt:'2026-09-24',updatedAt:'2026-09-24'};state.tags.push(t);return wait(t);},
    adminUpdateTag: async (id:number,input:any) => { const t=state.tags.find((x:any)=>x.tagId===Number(id));if(t)Object.assign(t,input);return wait(t);},
    adminGetTopics: async (q:any={}) => wait(page(state.topics,q)),
    adminCreateTopic: async (input:any) => { const t={topicId:++nextId,name:input.name,description:input.description||'',coverImage:'',postCount:0,participantCount:0,isOfficial:false,isRecommended:false,status:'enabled',createdAt:'2026-09-24',updatedAt:'2026-09-24'};state.topics.push(t);return wait(t);},
    adminUpdateTopic: async (id:number,input:any) => { const t=state.topics.find((x:any)=>x.topicId===Number(id));if(t)Object.assign(t,input);return wait(t);},
    adminGetCircles: async (q:any={}) => wait(page(state.circles,q)),
    adminStats: async () => wait({userCount:state.users.length,postCount:state.posts.length,commentCount:state.comments.length,circleCount:state.circles.length}),
    getAdminLogs: async (q:any={}) => wait(page([{logId:1,adminName:'Lin Chen',action:'client-platform-refactor',targetType:'repository',targetId:1,detail:'启用 Feedora Client Platform Mock 模式',createdAt:'2026-09-24 12:00:00'}],q)),
    getBaseModels: async (q:any={}) => wait(page(state.baseModels,q)),
    saveBaseModel: async (input:any) => { let m=state.baseModels.find((x:any)=>x.modelId===Number(input.modelId));if(m)Object.assign(m,input);else{m={modelId:++nextId,...input,createdAt:'2026-09-24',updatedAt:'2026-09-24'};state.baseModels.push(m);}return wait(m);},
    setDefaultModel: async (id:number) => { state.baseModels.forEach((x:any)=>x.isDefault=x.modelId===Number(id));return wait({success:true});},
    setModelStatus: async (id:number,status:string) => { const m=state.baseModels.find((x:any)=>x.modelId===Number(id));if(m)m.status=status;return wait(m);},
    getTokenConfig: async () => wait(state.tokenConfig),
    saveTokenConfig: async (input:any) => {Object.assign(state.tokenConfig,input,{updatedAt:'2026-09-24'});return wait(state.tokenConfig);},
    getTokenUsages: async (q:any={}) => wait(page(state.users.map((u:any,i:number)=>({userId:u.userId,nickname:u.nickname,avatar:u.avatar,todayTokens:18000-i*2500,monthTokens:360000-i*48000,lastUsedAt:'2026-09-24 11:58:00'})),q)),
    resetUserTokens: async () => wait({success:true}),
  };

  return new Proxy(methods,{
    get(target,prop){
      const name=String(prop);
      if(name in target) return target[name];
      return async (...args:any[]) => {
        console.warn('[mock-api] unimplemented method:',name,args);
        return wait(undefined);
      };
    }
  });
}

export const mockApi = createMockFeedoraApi();
