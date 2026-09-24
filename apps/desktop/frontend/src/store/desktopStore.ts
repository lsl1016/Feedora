import { create } from 'zustand';
import {
  seedKnowledgeBases,
  seedNotes,
  seedNotifications,
  seedPosts,
  seedRepositories,
  seedRepositoryFiles,
  seedSearchResults,
} from '@feedora/mock-data';

export interface ContextItem {
  id: string;
  type: 'post' | 'note' | 'knowledge-base' | 'repository' | 'file' | 'search';
  label: string;
}

interface DesktopSettings {
  serverUrl: string;
  runtimeUrl: string;
  workspacePath: string;
  importPath: string;
  indexMode: 'local' | 'hybrid';
  incrementalIndex: boolean;
  autoIndex: boolean;
  allowAiLocalFiles: boolean;
  encryptLocalData: boolean;
}

interface DesktopState {
  posts: any[];
  notes: any[];
  knowledgeBases: any[];
  notifications: any[];
  repositories: any[];
  repositoryFiles: any[];
  searchResults: any[];
  context: ContextItem[];
  aiMessages: Array<{ id: number; role: 'user' | 'assistant'; content: string }>;
  selectedRepositoryId: string;
  settings: DesktopSettings;
  toggleLike: (postId: number) => void;
  toggleFavorite: (postId: number) => void;
  savePostAsNote: (postId: number) => number | undefined;
  saveNote: (note: any) => number;
  createKnowledgeBase: (name: string, description?: string) => number;
  publishNote: (noteId: number) => number | undefined;
  markNotificationRead: (notificationId: number) => void;
  markAllNotificationsRead: () => void;
  addContext: (item: ContextItem) => void;
  removeContext: (id: string) => void;
  clearContext: () => void;
  askAi: (text: string) => void;
  saveLastAiAsNote: () => number | undefined;
  selectRepository: (id: string) => void;
  updateSettings: (patch: Partial<DesktopSettings>) => void;
}

const copy = <T,>(value: T): T => JSON.parse(JSON.stringify(value));
let sequence = 9000;

function buildAssistantReply(text: string, context: ContextItem[]) {
  const scope = context.length
    ? '当前已使用 ' + context.map((item) => item.label).join('、') + ' 作为上下文。'
    : '当前没有附加上下文，我会基于页面 Mock 数据回答。';

  if (/周报|总结/.test(text)) {
    return scope + '\n\n本次可整理为三部分：\n1. Client Platform 已完成 Web/Desktop 双宿主拆分；\n2. Desktop 已具备 Workspace、Repository、Search、AI Panel 的产品骨架；\n3. 下一步应补齐 SQLite/FTS5 与真实 Agent Runtime Adapter。';
  }
  if (/发布|帖子/.test(text)) {
    return scope + '\n\n建议发布结构：背景与问题 → 架构方案 → 模块边界 → 关键交互 → 当前实现 → 后续规划。发布前应检查是否包含私有代码、Token 或内部地址。';
  }
  if (/仓库|文件|代码/.test(text)) {
    return scope + '\n\n从客户端边界看，这部分代码应该保持“React Feature → Port → Wails Adapter → Go Service”的依赖方向，避免页面直接调用 Wails Binding。';
  }
  return scope + '\n\n可以继续围绕产品边界、数据流、模块职责和交互状态展开。我建议把结论保存为笔记，再决定是否发布到社区。';
}

export const useDesktopStore = create<DesktopState>((set, get) => ({
  posts: copy(seedPosts),
  notes: copy(seedNotes),
  knowledgeBases: copy(seedKnowledgeBases),
  notifications: copy(seedNotifications),
  repositories: copy(seedRepositories),
  repositoryFiles: copy(seedRepositoryFiles),
  searchResults: copy(seedSearchResults),
  context: [],
  aiMessages: [
    { id: 1, role: 'assistant', content: 'Feedora AI 助手已就绪。你可以把帖子、笔记、仓库文件或搜索结果加入 Context。' },
  ],
  selectedRepositoryId: seedRepositories[0]?.id || '',
  settings: {
    serverUrl: 'https://cloud.feedora.com',
    runtimeUrl: 'http://localhost:8080',
    workspacePath: 'D:/Feedora/Workspace',
    importPath: 'D:/Feedora/Inbox',
    indexMode: 'hybrid',
    incrementalIndex: true,
    autoIndex: true,
    allowAiLocalFiles: true,
    encryptLocalData: true,
  },

  toggleLike: (postId) => set((state) => ({
    posts: state.posts.map((post) => post.postId === postId
      ? { ...post, liked: !post.liked, likeCount: post.likeCount + (post.liked ? -1 : 1) }
      : post),
  })),

  toggleFavorite: (postId) => set((state) => ({
    posts: state.posts.map((post) => post.postId === postId
      ? { ...post, favorited: !post.favorited, favoriteCount: post.favoriteCount + (post.favorited ? -1 : 1) }
      : post),
  })),

  savePostAsNote: (postId) => {
    const post = get().posts.find((item) => item.postId === postId);
    if (!post) return undefined;
    const id = ++sequence;
    set((state) => ({
      notes: [{
        noteId: id,
        title: post.title,
        content: '# ' + post.title + '\n\n' + post.content + '\n\n> 来源：Feedora Community /posts/' + post.postId,
        summary: post.summary,
        knowledgeBaseId: 1,
        knowledgeBaseName: 'Feedora Desktop',
        ownerId: 1,
        status: 'normal',
        createdAt: '刚刚',
        updatedAt: '刚刚',
      }, ...state.notes],
    }));
    return id;
  },

  saveNote: (note) => {
    const existingId = Number(note.noteId || 0);
    if (existingId) {
      set((state) => ({
        notes: state.notes.map((item) => item.noteId === existingId
          ? { ...item, ...note, summary: String(note.content || '').replace(/[#>*_`]/g, '').slice(0, 100), updatedAt: '刚刚' }
          : item),
      }));
      return existingId;
    }
    const id = ++sequence;
    set((state) => ({
      notes: [{
        noteId: id,
        title: note.title || '未命名笔记',
        content: note.content || '',
        summary: String(note.content || '').replace(/[#>*_`]/g, '').slice(0, 100),
        knowledgeBaseId: note.knowledgeBaseId || 1,
        knowledgeBaseName: state.knowledgeBases.find((kb) => kb.knowledgeBaseId === (note.knowledgeBaseId || 1))?.name,
        ownerId: 1,
        status: 'normal',
        createdAt: '刚刚',
        updatedAt: '刚刚',
      }, ...state.notes],
    }));
    return id;
  },

  createKnowledgeBase: (name, description = '') => {
    const id = ++sequence;
    set((state) => ({
      knowledgeBases: [{
        knowledgeBaseId: id,
        name,
        description,
        ownerId: 1,
        noteCount: 0,
        memberCount: 1,
        visibility: 'private',
        createdAt: '刚刚',
        updatedAt: '刚刚',
      }, ...state.knowledgeBases],
    }));
    return id;
  },

  publishNote: (noteId) => {
    const note = get().notes.find((item) => item.noteId === noteId);
    if (!note) return undefined;
    const id = ++sequence;
    set((state) => ({
      posts: [{
        postId: id,
        postType: 'original',
        authorId: 1,
        author: state.posts[0]?.author,
        title: note.title,
        content: note.content,
        summary: note.summary,
        images: [],
        tags: [{ tagId: 5, tagName: '架构设计' }],
        topics: [],
        visibility: 'public',
        status: 'published',
        isTop: false,
        isFeatured: false,
        isSelected: false,
        viewCount: 0,
        likeCount: 0,
        commentCount: 0,
        favoriteCount: 0,
        shareCount: 0,
        repostCount: 0,
        hotScore: 0,
        liked: false,
        favorited: false,
        followedAuthor: false,
        createdAt: '刚刚',
        publishedAt: '刚刚',
        updatedAt: '刚刚',
      }, ...state.posts],
    }));
    return id;
  },

  markNotificationRead: (notificationId) => set((state) => ({
    notifications: state.notifications.map((item) => item.notificationId === notificationId
      ? { ...item, readStatus: 'read' }
      : item),
  })),

  markAllNotificationsRead: () => set((state) => ({
    notifications: state.notifications.map((item) => ({ ...item, readStatus: 'read' })),
  })),

  addContext: (item) => set((state) => state.context.some((entry) => entry.id === item.id)
    ? state
    : { context: [...state.context, item] }),

  removeContext: (id) => set((state) => ({ context: state.context.filter((item) => item.id !== id) })),
  clearContext: () => set({ context: [] }),

  askAi: (text) => {
    if (!text.trim()) return;
    const state = get();
    const user = { id: ++sequence, role: 'user' as const, content: text.trim() };
    const assistant = { id: ++sequence, role: 'assistant' as const, content: buildAssistantReply(text, state.context) };
    set((current) => ({ aiMessages: [...current.aiMessages, user, assistant] }));
  },

  saveLastAiAsNote: () => {
    const message = [...get().aiMessages].reverse().find((item) => item.role === 'assistant');
    if (!message) return undefined;
    return get().saveNote({
      title: 'AI 结论 ' + new Date().toLocaleTimeString(),
      content: message.content,
      knowledgeBaseId: 1,
    });
  },

  selectRepository: (id) => set({ selectedRepositoryId: id }),
  updateSettings: (patch) => set((state) => ({ settings: { ...state.settings, ...patch } })),
}));
