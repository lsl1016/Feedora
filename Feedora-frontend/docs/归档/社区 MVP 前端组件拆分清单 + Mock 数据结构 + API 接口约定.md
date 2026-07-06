# 社区 MVP 前端组件拆分清单 + Mock 数据结构 + API 接口约定

## 1. 推荐前端技术栈

推荐使用：

```text
React + TypeScript + React Router + Zustand + Axios + Ant Design
```

如果你更偏 Vue，也可以替换为：

```text
Vue 3 + TypeScript + Vue Router + Pinia + Axios + Element Plus
```

本文档默认以 **React + TypeScript + Ant Design** 为描述对象。

---

# 2. 前端目录结构建议

```text
src
├── api
│   ├── request.ts
│   ├── auth.ts
│   ├── user.ts
│   ├── post.ts
│   ├── comment.ts
│   ├── interaction.ts
│   ├── follow.ts
│   ├── notification.ts
│   ├── report.ts
│   ├── admin.ts
│   └── mock.ts
│
├── assets
│   ├── images
│   └── icons
│
├── components
│   ├── common
│   ├── layout
│   ├── post
│   ├── comment
│   ├── user
│   ├── notification
│   └── admin
│
├── constants
│   ├── enums.ts
│   ├── route.ts
│   └── options.ts
│
├── hooks
│   ├── useAuth.ts
│   ├── usePagination.ts
│   ├── usePostActions.ts
│   └── useAdminTable.ts
│
├── layouts
│   ├── FrontLayout.tsx
│   └── AdminLayout.tsx
│
├── pages
│   ├── auth
│   ├── home
│   ├── hot
│   ├── following
│   ├── post
│   ├── user
│   ├── notification
│   └── admin
│
├── router
│   └── index.tsx
│
├── store
│   ├── authStore.ts
│   ├── userStore.ts
│   └── appStore.ts
│
├── styles
│   ├── global.css
│   └── variables.css
│
├── types
│   ├── common.ts
│   ├── user.ts
│   ├── post.ts
│   ├── comment.ts
│   ├── notification.ts
│   ├── report.ts
│   └── admin.ts
│
├── utils
│   ├── format.ts
│   ├── storage.ts
│   ├── permission.ts
│   └── validator.ts
│
├── App.tsx
└── main.tsx
```

---

# 3. 页面级组件拆分

## 3.1 前台页面

```text
pages
├── auth
│   ├── LoginPage.tsx
│   └── RegisterPage.tsx
│
├── home
│   └── HomePage.tsx
│
├── hot
│   └── HotPage.tsx
│
├── following
│   └── FollowingPage.tsx
│
├── post
│   ├── PostDetailPage.tsx
│   ├── PostCreatePage.tsx
│   └── PostEditPage.tsx
│
├── user
│   ├── MyProfilePage.tsx
│   └── UserProfilePage.tsx
│
└── notification
    └── NotificationPage.tsx
```

### 页面职责说明

| 页面    | 组件                 | 说明               |
| ----- | ------------------ | ---------------- |
| 登录页   | `LoginPage`        | 用户登录             |
| 注册页   | `RegisterPage`     | 用户注册             |
| 首页    | `HomePage`         | 全站帖子流            |
| 热门页   | `HotPage`          | 今日热门、本周热门、总榜     |
| 关注页   | `FollowingPage`    | 关注用户发布的帖子        |
| 帖子详情页 | `PostDetailPage`   | 帖子正文、互动、评论       |
| 发布帖子页 | `PostCreatePage`   | 新建帖子             |
| 编辑帖子页 | `PostEditPage`     | 修改已有帖子           |
| 我的主页  | `MyProfilePage`    | 我的资料、帖子、点赞、收藏、关注 |
| 用户主页  | `UserProfilePage`  | 其他用户主页           |
| 通知页   | `NotificationPage` | 通知中心             |

---

## 3.2 后台页面

```text
pages
└── admin
    ├── AdminDashboardPage.tsx
    ├── AdminUserPage.tsx
    ├── AdminPostPage.tsx
    ├── AdminCommentPage.tsx
    ├── AdminReviewPage.tsx
    ├── AdminReportPage.tsx
    ├── AdminSensitiveWordPage.tsx
    └── AdminOperationLogPage.tsx
```

### 页面职责说明

| 页面    | 组件                       | 说明            |
| ----- | ------------------------ | ------------- |
| 数据看板  | `AdminDashboardPage`     | 社区核心指标与趋势     |
| 用户管理  | `AdminUserPage`          | 用户查询、禁言、封禁、解封 |
| 帖子管理  | `AdminPostPage`          | 帖子查询、审核、下架、删除 |
| 评论管理  | `AdminCommentPage`       | 评论查询、查看上下文、删除 |
| 审核管理  | `AdminReviewPage`        | 帖子审核、评论审核     |
| 举报管理  | `AdminReportPage`        | 处理举报          |
| 敏感词管理 | `AdminSensitiveWordPage` | 敏感词增删改查       |
| 操作日志  | `AdminOperationLogPage`  | 后台操作审计        |

---

# 4. 布局组件拆分

## 4.1 前台布局组件

```text
layouts
└── FrontLayout.tsx
```

## FrontLayout 结构

```text
FrontLayout
├── FrontHeader
│   ├── Logo
│   ├── NavMenu
│   ├── SearchBox
│   ├── PublishButton
│   ├── NotificationEntry
│   └── UserAvatarDropdown
└── MainContent
```

### 组件说明

| 组件                   | 说明            |
| -------------------- | ------------- |
| `FrontHeader`        | 前台顶部导航        |
| `Logo`               | 社区名称，点击跳转首页   |
| `NavMenu`            | 首页、热门、关注      |
| `SearchBox`          | 一期可展示，不实现复杂搜索 |
| `PublishButton`      | 发布帖子按钮        |
| `NotificationEntry`  | 通知入口，展示未读红点   |
| `UserAvatarDropdown` | 用户头像下拉菜单      |
| `MainContent`        | 页面主体区域        |

---

## 4.2 后台布局组件

```text
layouts
└── AdminLayout.tsx
```

## AdminLayout 结构

```text
AdminLayout
├── AdminSidebar
│   ├── DashboardMenuItem
│   ├── UserMenuItem
│   ├── PostMenuItem
│   ├── CommentMenuItem
│   ├── ReviewMenuItem
│   ├── ReportMenuItem
│   ├── SensitiveWordMenuItem
│   └── OperationLogMenuItem
├── AdminTopbar
└── AdminContent
```

### 组件说明

| 组件                  | 说明      |
| ------------------- | ------- |
| `AdminSidebar`      | 后台左侧菜单  |
| `AdminTopbar`       | 后台顶部栏   |
| `AdminContent`      | 后台内容区域  |
| `AdminBreadcrumb`   | 面包屑，可选  |
| `AdminUserDropdown` | 管理员信息菜单 |

---

# 5. 通用组件拆分

## 5.1 common 通用组件

```text
components
└── common
    ├── PageContainer.tsx
    ├── CardContainer.tsx
    ├── EmptyState.tsx
    ├── ErrorState.tsx
    ├── LoadingSkeleton.tsx
    ├── ConfirmModal.tsx
    ├── StatusTag.tsx
    ├── ImagePreviewModal.tsx
    ├── UploadImageGrid.tsx
    ├── FormErrorText.tsx
    └── PaginationBar.tsx
```

### 组件说明

| 组件                  | 说明       |
| ------------------- | -------- |
| `PageContainer`     | 页面统一外层容器 |
| `CardContainer`     | 白色卡片容器   |
| `EmptyState`        | 空状态      |
| `ErrorState`        | 错误状态     |
| `LoadingSkeleton`   | 骨架屏      |
| `ConfirmModal`      | 通用确认弹窗   |
| `StatusTag`         | 状态标签     |
| `ImagePreviewModal` | 图片预览弹窗   |
| `UploadImageGrid`   | 图片上传宫格   |
| `FormErrorText`     | 表单错误提示   |
| `PaginationBar`     | 分页组件     |

---

# 6. 业务组件拆分

## 6.1 帖子组件

```text
components
└── post
    ├── PostCard.tsx
    ├── PostList.tsx
    ├── PostFeedTabs.tsx
    ├── PublishEntryCard.tsx
    ├── PostActionBar.tsx
    ├── PostAuthorInfo.tsx
    ├── PostContent.tsx
    ├── PostForm.tsx
    ├── PostStatusTag.tsx
    ├── HotRankList.tsx
    └── HotRankItem.tsx
```

### 组件说明

| 组件                 | 说明                             |
| ------------------ | ------------------------------ |
| `PostCard`         | 首页、关注流中的帖子卡片                   |
| `PostList`         | 帖子列表，负责 loading、empty、error、分页 |
| `PostFeedTabs`     | 最新 / 热门 Tab                    |
| `PublishEntryCard` | 首页顶部发布入口                       |
| `PostActionBar`    | 点赞、收藏、评论、举报操作栏                 |
| `PostAuthorInfo`   | 作者信息区                          |
| `PostContent`      | 帖子正文展示                         |
| `PostForm`         | 发布 / 编辑帖子表单                    |
| `PostStatusTag`    | 我的帖子状态标签                       |
| `HotRankList`      | 热门榜单列表                         |
| `HotRankItem`      | 热门榜单单项                         |

---

## 6.2 评论组件

```text
components
└── comment
    ├── CommentInput.tsx
    ├── CommentList.tsx
    ├── CommentItem.tsx
    ├── ReplyItem.tsx
    ├── ReplyInput.tsx
    └── CommentActionBar.tsx
```

### 组件说明

| 组件                 | 说明         |
| ------------------ | ---------- |
| `CommentInput`     | 帖子详情页评论输入框 |
| `CommentList`      | 评论列表       |
| `CommentItem`      | 一级评论       |
| `ReplyItem`        | 二级回复       |
| `ReplyInput`       | 回复输入框      |
| `CommentActionBar` | 评论点赞、回复、删除 |

---

## 6.3 用户组件

```text
components
└── user
    ├── UserProfileCard.tsx
    ├── UserStatBar.tsx
    ├── EditProfileModal.tsx
    ├── FollowButton.tsx
    ├── UserList.tsx
    ├── UserListItem.tsx
    └── BlockUserModal.tsx
```

### 组件说明

| 组件                 | 说明          |
| ------------------ | ----------- |
| `UserProfileCard`  | 用户资料卡片      |
| `UserStatBar`      | 关注数、粉丝数、获赞数 |
| `EditProfileModal` | 编辑资料弹窗      |
| `FollowButton`     | 关注 / 已关注按钮  |
| `UserList`         | 用户列表        |
| `UserListItem`     | 用户列表单项      |
| `BlockUserModal`   | 拉黑用户确认弹窗    |

---

## 6.4 通知组件

```text
components
└── notification
    ├── NotificationTabs.tsx
    ├── NotificationList.tsx
    ├── NotificationItem.tsx
    └── NotificationDetailModal.tsx
```

### 组件说明

| 组件                        | 说明       |
| ------------------------- | -------- |
| `NotificationTabs`        | 通知类型 Tab |
| `NotificationList`        | 通知列表     |
| `NotificationItem`        | 单条通知     |
| `NotificationDetailModal` | 系统通知详情弹窗 |

---

## 6.5 举报组件

```text
components
└── report
    └── ReportModal.tsx
```

### 组件说明

| 组件            | 说明              |
| ------------- | --------------- |
| `ReportModal` | 举报帖子、评论、用户的统一弹窗 |

---

## 6.6 后台组件

```text
components
└── admin
    ├── AdminMetricCard.tsx
    ├── AdminTrendChart.tsx
    ├── AdminFilterBar.tsx
    ├── AdminTableActions.tsx
    ├── UserDetailDrawer.tsx
    ├── PostPreviewDrawer.tsx
    ├── CommentContextDrawer.tsx
    ├── ReviewDetailDrawer.tsx
    ├── ReportTargetDrawer.tsx
    ├── RejectReasonModal.tsx
    ├── TakedownReasonModal.tsx
    ├── ReportHandleModal.tsx
    └── SensitiveWordFormModal.tsx
```

### 组件说明

| 组件                       | 说明           |
| ------------------------ | ------------ |
| `AdminMetricCard`        | 后台指标卡片       |
| `AdminTrendChart`        | 趋势图          |
| `AdminFilterBar`         | 筛选区          |
| `AdminTableActions`      | 表格操作按钮组      |
| `UserDetailDrawer`       | 用户详情抽屉       |
| `PostPreviewDrawer`      | 帖子预览抽屉       |
| `CommentContextDrawer`   | 评论上下文抽屉      |
| `ReviewDetailDrawer`     | 审核详情抽屉       |
| `ReportTargetDrawer`     | 举报对象详情抽屉     |
| `RejectReasonModal`      | 拒绝原因弹窗       |
| `TakedownReasonModal`    | 下架原因弹窗       |
| `ReportHandleModal`      | 举报处理弹窗       |
| `SensitiveWordFormModal` | 敏感词新增 / 编辑弹窗 |

---

# 7. TypeScript 类型定义

## 7.1 通用类型

```ts
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export interface PageRequest {
  page: number;
  pageSize: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  pageSize: number;
}

export type ID = number;

export type TargetType = 'post' | 'comment' | 'user';

export type SortType = 'latest' | 'hot';
```

---

## 7.2 用户类型

```ts
export type UserStatus = 'normal' | 'muted' | 'banned' | 'deleted';

export interface User {
  userId: number;
  account: string;
  nickname: string;
  avatar: string;
  bio: string;
  status: UserStatus;
  postCount: number;
  commentCount: number;
  followerCount: number;
  followingCount: number;
  likeReceivedCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface LoginRequest {
  account: string;
  password: string;
}

export interface RegisterRequest {
  account: string;
  nickname: string;
  password: string;
  confirmPassword: string;
}

export interface LoginResult {
  token: string;
  user: User;
}

export interface UpdateProfileRequest {
  nickname: string;
  avatar: string;
  bio: string;
}
```

---

## 7.3 帖子类型

```ts
export type PostStatus =
  | 'draft'
  | 'reviewing'
  | 'published'
  | 'rejected'
  | 'deleted'
  | 'takedown';

export interface Post {
  postId: number;
  authorId: number;
  author: UserSummary;
  title: string;
  content: string;
  summary: string;
  images: string[];
  status: PostStatus;
  viewCount: number;
  likeCount: number;
  commentCount: number;
  favoriteCount: number;
  hotScore: number;
  liked: boolean;
  favorited: boolean;
  followedAuthor: boolean;
  rejectReason?: string;
  takedownReason?: string;
  createdAt: string;
  publishedAt?: string;
  updatedAt: string;
}

export interface UserSummary {
  userId: number;
  nickname: string;
  avatar: string;
  bio?: string;
}

export interface CreatePostRequest {
  title: string;
  content: string;
  images: string[];
  saveAsDraft: boolean;
}

export interface UpdatePostRequest {
  title: string;
  content: string;
  images: string[];
}

export interface PostQuery {
  page: number;
  pageSize: number;
  sort?: SortType;
  status?: PostStatus;
  keyword?: string;
  authorId?: number;
}
```

---

## 7.4 评论类型

```ts
export type CommentStatus = 'normal' | 'deleted' | 'rejected';

export interface Comment {
  commentId: number;
  postId: number;
  userId: number;
  user: UserSummary;
  parentId?: number;
  rootId?: number;
  replyToUserId?: number;
  replyToUser?: UserSummary;
  content: string;
  likeCount: number;
  liked: boolean;
  status: CommentStatus;
  replies: Comment[];
  createdAt: string;
  updatedAt: string;
}

export interface CreateCommentRequest {
  postId: number;
  content: string;
}

export interface CreateReplyRequest {
  postId: number;
  parentId: number;
  rootId: number;
  replyToUserId: number;
  content: string;
}
```

---

## 7.5 通知类型

```ts
export type NotificationType =
  | 'like_post'
  | 'like_comment'
  | 'comment_post'
  | 'reply_comment'
  | 'follow_user'
  | 'review_pass'
  | 'review_reject'
  | 'system';

export type ReadStatus = 'unread' | 'read';

export interface Notification {
  notificationId: number;
  receiverId: number;
  senderId?: number;
  sender?: UserSummary;
  type: NotificationType;
  title: string;
  content: string;
  targetType?: TargetType;
  targetId?: number;
  extraData?: Record<string, unknown>;
  readStatus: ReadStatus;
  createdAt: string;
}
```

---

## 7.6 举报类型

```ts
export type ReportReason =
  | 'spam'
  | 'illegal'
  | 'abuse'
  | 'porn'
  | 'violation'
  | 'other';

export type ReportStatus = 'pending' | 'approved' | 'rejected' | 'closed';

export interface Report {
  reportId: number;
  reporterId: number;
  reporter: UserSummary;
  targetType: TargetType;
  targetId: number;
  reason: ReportReason;
  description: string;
  status: ReportStatus;
  handlerId?: number;
  handleResult?: string;
  handledAt?: string;
  createdAt: string;
}

export interface CreateReportRequest {
  targetType: TargetType;
  targetId: number;
  reason: ReportReason;
  description: string;
}
```

---

## 7.7 后台类型

```ts
export interface AdminDashboardStats {
  todayNewUserCount: number;
  todayPostCount: number;
  todayCommentCount: number;
  todayLikeCount: number;
  pendingReviewCount: number;
  pendingReportCount: number;
}

export interface TrendPoint {
  date: string;
  value: number;
}

export type SensitiveWordStrategy = 'block' | 'review';
export type EnabledStatus = 'enabled' | 'disabled';

export interface SensitiveWord {
  wordId: number;
  word: string;
  strategy: SensitiveWordStrategy;
  status: EnabledStatus;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface AdminOperationLog {
  logId: number;
  adminId: number;
  adminName: string;
  action: string;
  targetType: string;
  targetId: number;
  detail: string;
  ip: string;
  createdAt: string;
}
```

---

# 8. Mock 数据结构

## 8.1 Mock 用户

```ts
export const mockUsers: User[] = [
  {
    userId: 10001,
    account: 'zhangsan',
    nickname: '张三',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan',
    bio: '后端开发，喜欢 Go、Java 和数据平台。',
    status: 'normal',
    postCount: 12,
    commentCount: 35,
    followerCount: 256,
    followingCount: 86,
    likeReceivedCount: 1234,
    createdAt: '2026-07-01 10:00:00',
    updatedAt: '2026-07-05 10:00:00',
  },
  {
    userId: 10002,
    account: 'frontend_dev',
    nickname: '前端小李',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=frontend',
    bio: 'React / Vue 前端开发者。',
    status: 'normal',
    postCount: 8,
    commentCount: 22,
    followerCount: 128,
    followingCount: 45,
    likeReceivedCount: 560,
    createdAt: '2026-07-02 11:20:00',
    updatedAt: '2026-07-05 12:00:00',
  },
  {
    userId: 10003,
    account: 'admin',
    nickname: '社区管理员',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin',
    bio: '社区运营与内容治理。',
    status: 'normal',
    postCount: 3,
    commentCount: 6,
    followerCount: 500,
    followingCount: 10,
    likeReceivedCount: 2000,
    createdAt: '2026-06-30 09:00:00',
    updatedAt: '2026-07-05 09:00:00',
  },
];
```

---

## 8.2 Mock 帖子

```ts
export const mockPosts: Post[] = [
  {
    postId: 20001,
    authorId: 10002,
    author: {
      userId: 10002,
      nickname: '前端小李',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=frontend',
    },
    title: '如何高效学习一门新的编程语言？',
    content:
      '最近在学习 Go 语言，感觉语法并不难，难点在于工程实践。我的方法是先快速过一遍语法，然后直接做一个小项目，在项目中遇到问题再回头补知识点。',
    summary: '最近在学习 Go 语言，感觉语法并不难，难点在于工程实践...',
    images: [
      'https://images.unsplash.com/photo-1516321318423-f06f85e504b3',
      'https://images.unsplash.com/photo-1515879218367-8466d910aaa4',
    ],
    status: 'published',
    viewCount: 1280,
    likeCount: 432,
    commentCount: 32,
    favoriteCount: 88,
    hotScore: 2600,
    liked: false,
    favorited: false,
    followedAuthor: false,
    createdAt: '2026-07-05 09:30:00',
    publishedAt: '2026-07-05 09:35:00',
    updatedAt: '2026-07-05 09:35:00',
  },
  {
    postId: 20002,
    authorId: 10001,
    author: {
      userId: 10001,
      nickname: '张三',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan',
    },
    title: '三层架构风格的界面设计分享',
    content:
      '最近在设计一个社区项目，前端页面采用卡片式布局，后台使用管理系统结构。整体目标是让页面结构简单、清晰，便于后续接入真实 API。',
    summary: '最近在设计一个社区项目，前端页面采用卡片式布局...',
    images: [],
    status: 'published',
    viewCount: 980,
    likeCount: 156,
    commentCount: 24,
    favoriteCount: 45,
    hotScore: 1380,
    liked: true,
    favorited: false,
    followedAuthor: true,
    createdAt: '2026-07-04 18:00:00',
    publishedAt: '2026-07-04 18:05:00',
    updatedAt: '2026-07-04 18:05:00',
  },
  {
    postId: 20003,
    authorId: 10001,
    author: {
      userId: 10001,
      nickname: '张三',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan',
    },
    title: 'React 状态管理方案对比',
    content: '本文简单对比 Zustand、Redux Toolkit 和 Context 在中小型项目中的使用体验。',
    summary: '本文简单对比 Zustand、Redux Toolkit 和 Context...',
    images: [],
    status: 'reviewing',
    viewCount: 0,
    likeCount: 0,
    commentCount: 0,
    favoriteCount: 0,
    hotScore: 0,
    liked: false,
    favorited: false,
    followedAuthor: false,
    createdAt: '2026-07-05 14:00:00',
    updatedAt: '2026-07-05 14:00:00',
  },
];
```

---

## 8.3 Mock 评论

```ts
export const mockComments: Comment[] = [
  {
    commentId: 30001,
    postId: 20001,
    userId: 10001,
    user: {
      userId: 10001,
      nickname: '张三',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan',
    },
    content: '这个学习方法挺实用的，先做项目再补知识点确实效率高。',
    likeCount: 12,
    liked: false,
    status: 'normal',
    replies: [
      {
        commentId: 30002,
        postId: 20001,
        userId: 10002,
        user: {
          userId: 10002,
          nickname: '前端小李',
          avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=frontend',
        },
        parentId: 30001,
        rootId: 30001,
        replyToUserId: 10001,
        replyToUser: {
          userId: 10001,
          nickname: '张三',
          avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan',
        },
        content: '是的，尤其是做完整项目时，会暴露很多真实问题。',
        likeCount: 3,
        liked: false,
        status: 'normal',
        replies: [],
        createdAt: '2026-07-05 10:20:00',
        updatedAt: '2026-07-05 10:20:00',
      },
    ],
    createdAt: '2026-07-05 10:10:00',
    updatedAt: '2026-07-05 10:10:00',
  },
];
```

---

## 8.4 Mock 通知

```ts
export const mockNotifications: Notification[] = [
  {
    notificationId: 40001,
    receiverId: 10001,
    senderId: 10002,
    sender: {
      userId: 10002,
      nickname: '前端小李',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=frontend',
    },
    type: 'comment_post',
    title: '你的帖子收到了新的评论',
    content: '前端小李 评论了你的帖子《三层架构风格的界面设计分享》',
    targetType: 'post',
    targetId: 20002,
    extraData: {
      commentId: 30001,
    },
    readStatus: 'unread',
    createdAt: '2026-07-05 10:30:00',
  },
  {
    notificationId: 40002,
    receiverId: 10001,
    senderId: 10002,
    sender: {
      userId: 10002,
      nickname: '前端小李',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=frontend',
    },
    type: 'follow_user',
    title: '你有新的关注者',
    content: '前端小李 关注了你',
    targetType: 'user',
    targetId: 10002,
    readStatus: 'read',
    createdAt: '2026-07-04 18:20:00',
  },
];
```

---

## 8.5 Mock 举报

```ts
export const mockReports: Report[] = [
  {
    reportId: 50001,
    reporterId: 10001,
    reporter: {
      userId: 10001,
      nickname: '张三',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan',
    },
    targetType: 'post',
    targetId: 20001,
    reason: 'spam',
    description: '疑似广告内容。',
    status: 'pending',
    createdAt: '2026-07-05 11:00:00',
  },
];
```

---

## 8.6 Mock 敏感词

```ts
export const mockSensitiveWords: SensitiveWord[] = [
  {
    wordId: 60001,
    word: '广告引流',
    strategy: 'review',
    status: 'enabled',
    createdBy: 10003,
    createdAt: '2026-07-01 10:00:00',
    updatedAt: '2026-07-01 10:00:00',
  },
  {
    wordId: 60002,
    word: '违法内容',
    strategy: 'block',
    status: 'enabled',
    createdBy: 10003,
    createdAt: '2026-07-01 10:10:00',
    updatedAt: '2026-07-01 10:10:00',
  },
];
```

---

## 8.7 Mock 后台数据

```ts
export const mockAdminDashboardStats: AdminDashboardStats = {
  todayNewUserCount: 128,
  todayPostCount: 256,
  todayCommentCount: 1234,
  todayLikeCount: 3200,
  pendingReviewCount: 18,
  pendingReportCount: 32,
};

export const mockUserTrend: TrendPoint[] = [
  { date: '07-01', value: 20 },
  { date: '07-02', value: 35 },
  { date: '07-03', value: 42 },
  { date: '07-04', value: 60 },
  { date: '07-05', value: 88 },
];

export const mockPostTrend: TrendPoint[] = [
  { date: '07-01', value: 50 },
  { date: '07-02', value: 80 },
  { date: '07-03', value: 120 },
  { date: '07-04', value: 160 },
  { date: '07-05', value: 256 },
];
```

---

# 9. API 接口通用约定

## 9.1 基础地址

```text
/api/v1
```

## 9.2 通用响应结构

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 9.3 分页响应结构

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "pageSize": 10
  }
}
```

## 9.4 错误响应结构

```json
{
  "code": 40001,
  "message": "账号或密码错误",
  "data": null
}
```

## 9.5 常见错误码

|   错误码 | 说明          |
| ----: | ----------- |
|     0 | 成功          |
| 40000 | 请求参数错误      |
| 40001 | 账号或密码错误     |
| 40100 | 未登录         |
| 40300 | 无权限         |
| 40301 | 用户已被禁言      |
| 40302 | 用户已被封禁      |
| 40400 | 资源不存在       |
| 40900 | 数据冲突，例如重复点赞 |
| 50000 | 系统异常        |

## 9.6 请求头约定

登录后接口统一携带：

```text
Authorization: Bearer <token>
```

---

# 10. 认证接口

## 10.1 登录

```text
POST /api/v1/auth/login
```

### 请求体

```json
{
  "account": "zhangsan",
  "password": "123456"
}
```

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "mock-token",
    "user": {
      "userId": 10001,
      "account": "zhangsan",
      "nickname": "张三",
      "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan",
      "bio": "后端开发，喜欢 Go、Java 和数据平台。",
      "status": "normal",
      "postCount": 12,
      "commentCount": 35,
      "followerCount": 256,
      "followingCount": 86,
      "likeReceivedCount": 1234,
      "createdAt": "2026-07-01 10:00:00",
      "updatedAt": "2026-07-05 10:00:00"
    }
  }
}
```

---

## 10.2 注册

```text
POST /api/v1/auth/register
```

### 请求体

```json
{
  "account": "lisi",
  "nickname": "李四",
  "password": "123456",
  "confirmPassword": "123456"
}
```

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "mock-token",
    "user": {}
  }
}
```

---

## 10.3 获取当前登录用户

```text
GET /api/v1/auth/me
```

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

## 10.4 退出登录

```text
POST /api/v1/auth/logout
```

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": true
}
```

---

# 11. 用户接口

## 11.1 获取用户详情

```text
GET /api/v1/users/{userId}
```

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

## 11.2 更新当前用户资料

```text
PUT /api/v1/users/me/profile
```

### 请求体

```json
{
  "nickname": "新的昵称",
  "avatar": "https://xxx.com/avatar.png",
  "bio": "新的个人简介"
}
```

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

---

## 11.3 获取我的帖子

```text
GET /api/v1/users/me/posts
```

### Query

| 参数       | 类型     | 说明      |
| -------- | ------ | ------- |
| page     | number | 页码      |
| pageSize | number | 每页数量    |
| status   | string | 帖子状态，可选 |

---

## 11.4 获取我的点赞帖子

```text
GET /api/v1/users/me/liked-posts
```

---

## 11.5 获取我的收藏帖子

```text
GET /api/v1/users/me/favorite-posts
```

---

## 11.6 获取我的关注用户

```text
GET /api/v1/users/me/following
```

---

## 11.7 获取用户公开帖子

```text
GET /api/v1/users/{userId}/posts
```

---

# 12. 帖子接口

## 12.1 获取首页帖子流

```text
GET /api/v1/posts
```

### Query

| 参数       | 类型     | 必填 | 说明           |
| -------- | ------ | -: | ------------ |
| page     | number |  是 | 页码           |
| pageSize | number |  是 | 每页数量         |
| sort     | string |  否 | latest / hot |
| keyword  | string |  否 | 搜索关键词        |

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "pageSize": 10
  }
}
```

---

## 12.2 获取关注流

```text
GET /api/v1/posts/following
```

### Query

| 参数       | 类型     | 必填 |
| -------- | ------ | -: |
| page     | number |  是 |
| pageSize | number |  是 |

---

## 12.3 获取热门榜单

```text
GET /api/v1/posts/hot
```

### Query

| 参数    | 类型     | 说明                 |
| ----- | ------ | ------------------ |
| type  | string | today / week / all |
| limit | number | 榜单数量               |

---

## 12.4 获取帖子详情

```text
GET /api/v1/posts/{postId}
```

---

## 12.5 创建帖子

```text
POST /api/v1/posts
```

### 请求体

```json
{
  "title": "如何高效学习一门新的编程语言？",
  "content": "正文内容",
  "images": ["https://xxx.com/1.png"],
  "saveAsDraft": false
}
```

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "postId": 20001,
    "status": "published"
  }
}
```

---

## 12.6 更新帖子

```text
PUT /api/v1/posts/{postId}
```

### 请求体

```json
{
  "title": "新的标题",
  "content": "新的正文",
  "images": []
}
```

---

## 12.7 删除帖子

```text
DELETE /api/v1/posts/{postId}
```

---

# 13. 评论接口

## 13.1 获取帖子评论

```text
GET /api/v1/posts/{postId}/comments
```

### Query

| 参数       | 类型     | 必填 |
| -------- | ------ | -: |
| page     | number |  是 |
| pageSize | number |  是 |

---

## 13.2 发表评论

```text
POST /api/v1/comments
```

### 请求体

```json
{
  "postId": 20001,
  "content": "这个学习方法挺实用的。"
}
```

---

## 13.3 回复评论

```text
POST /api/v1/comments/{commentId}/replies
```

### 请求体

```json
{
  "postId": 20001,
  "rootId": 30001,
  "replyToUserId": 10001,
  "content": "确实如此。"
}
```

---

## 13.4 删除评论

```text
DELETE /api/v1/comments/{commentId}
```

---

# 14. 互动接口

## 14.1 点赞帖子 / 取消点赞帖子

```text
POST /api/v1/posts/{postId}/like
DELETE /api/v1/posts/{postId}/like
```

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "liked": true,
    "likeCount": 433
  }
}
```

---

## 14.2 收藏帖子 / 取消收藏帖子

```text
POST /api/v1/posts/{postId}/favorite
DELETE /api/v1/posts/{postId}/favorite
```

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "favorited": true,
    "favoriteCount": 89
  }
}
```

---

## 14.3 点赞评论 / 取消点赞评论

```text
POST /api/v1/comments/{commentId}/like
DELETE /api/v1/comments/{commentId}/like
```

---

# 15. 关注接口

## 15.1 关注用户

```text
POST /api/v1/users/{userId}/follow
```

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "followed": true,
    "followerCount": 257
  }
}
```

---

## 15.2 取消关注用户

```text
DELETE /api/v1/users/{userId}/follow
```

---

## 15.3 拉黑用户

```text
POST /api/v1/users/{userId}/block
```

---

# 16. 通知接口

## 16.1 获取通知列表

```text
GET /api/v1/notifications
```

### Query

| 参数       | 类型     | 说明                                                      |
| -------- | ------ | ------------------------------------------------------- |
| type     | string | all / like / comment / reply / follow / review / system |
| page     | number | 页码                                                      |
| pageSize | number | 每页数量                                                    |

---

## 16.2 标记单条通知已读

```text
PUT /api/v1/notifications/{notificationId}/read
```

---

## 16.3 全部通知已读

```text
PUT /api/v1/notifications/read-all
```

---

## 16.4 获取未读通知数

```text
GET /api/v1/notifications/unread-count
```

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "count": 12
  }
}
```

---

# 17. 举报接口

## 17.1 提交举报

```text
POST /api/v1/reports
```

### 请求体

```json
{
  "targetType": "post",
  "targetId": 20001,
  "reason": "spam",
  "description": "疑似广告内容"
}
```

---

# 18. 文件上传接口

## 18.1 上传图片

```text
POST /api/v1/upload/image
```

### Content-Type

```text
multipart/form-data
```

### FormData

| 字段   | 类型   | 说明   |
| ---- | ---- | ---- |
| file | File | 图片文件 |

### 响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "url": "https://xxx.com/image.png"
  }
}
```

---

# 19. 后台接口

## 19.1 后台数据看板

```text
GET /api/v1/admin/dashboard/stats
GET /api/v1/admin/dashboard/user-trend
GET /api/v1/admin/dashboard/post-trend
```

---

## 19.2 后台用户管理

### 用户列表

```text
GET /api/v1/admin/users
```

### Query

| 参数        | 类型     | 说明     |
| --------- | ------ | ------ |
| userId    | number | 用户 ID  |
| nickname  | string | 昵称     |
| status    | string | 用户状态   |
| startDate | string | 注册开始时间 |
| endDate   | string | 注册结束时间 |
| page      | number | 页码     |
| pageSize  | number | 每页数量   |

### 禁言用户

```text
PUT /api/v1/admin/users/{userId}/mute
```

### 封禁用户

```text
PUT /api/v1/admin/users/{userId}/ban
```

### 解封用户

```text
PUT /api/v1/admin/users/{userId}/unban
```

---

## 19.3 后台帖子管理

### 帖子列表

```text
GET /api/v1/admin/posts
```

### 审核通过

```text
PUT /api/v1/admin/posts/{postId}/approve
```

### 审核拒绝

```text
PUT /api/v1/admin/posts/{postId}/reject
```

请求体：

```json
{
  "reason": "内容不符合社区规范"
}
```

### 下架帖子

```text
PUT /api/v1/admin/posts/{postId}/takedown
```

请求体：

```json
{
  "reason": "内容存在违规信息"
}
```

### 删除帖子

```text
DELETE /api/v1/admin/posts/{postId}
```

---

## 19.4 后台评论管理

### 评论列表

```text
GET /api/v1/admin/comments
```

### 删除评论

```text
DELETE /api/v1/admin/comments/{commentId}
```

---

## 19.5 后台审核管理

### 待审核帖子

```text
GET /api/v1/admin/reviews/posts
```

### 待审核评论

```text
GET /api/v1/admin/reviews/comments
```

### 审核通过

```text
PUT /api/v1/admin/reviews/{targetType}/{targetId}/approve
```

### 审核拒绝

```text
PUT /api/v1/admin/reviews/{targetType}/{targetId}/reject
```

请求体：

```json
{
  "reason": "内容包含不合规信息"
}
```

---

## 19.6 后台举报管理

### 举报列表

```text
GET /api/v1/admin/reports
```

### 通过举报

```text
PUT /api/v1/admin/reports/{reportId}/approve
```

请求体：

```json
{
  "action": "takedown_post",
  "handleResult": "举报成立，已下架帖子"
}
```

action 可选值：

```text
takedown_post
delete_comment
mute_user
ban_user
```

### 驳回举报

```text
PUT /api/v1/admin/reports/{reportId}/reject
```

请求体：

```json
{
  "handleResult": "举报不成立"
}
```

---

## 19.7 后台敏感词管理

### 敏感词列表

```text
GET /api/v1/admin/sensitive-words
```

### 新增敏感词

```text
POST /api/v1/admin/sensitive-words
```

请求体：

```json
{
  "word": "广告引流",
  "strategy": "review",
  "status": "enabled"
}
```

### 编辑敏感词

```text
PUT /api/v1/admin/sensitive-words/{wordId}
```

### 删除敏感词

```text
DELETE /api/v1/admin/sensitive-words/{wordId}
```

### 启用敏感词

```text
PUT /api/v1/admin/sensitive-words/{wordId}/enable
```

### 停用敏感词

```text
PUT /api/v1/admin/sensitive-words/{wordId}/disable
```

---

## 19.8 后台操作日志

```text
GET /api/v1/admin/operation-logs
```

### Query

| 参数         | 类型     | 说明    |
| ---------- | ------ | ----- |
| adminName  | string | 管理员名称 |
| action     | string | 操作类型  |
| targetType | string | 对象类型  |
| startDate  | string | 开始时间  |
| endDate    | string | 结束时间  |
| page       | number | 页码    |
| pageSize   | number | 每页数量  |

---

# 20. 前端状态管理建议

## 20.1 authStore

```ts
interface AuthState {
  token: string | null;
  currentUser: User | null;
  isLogin: boolean;
  setToken: (token: string | null) => void;
  setCurrentUser: (user: User | null) => void;
  logout: () => void;
}
```

职责：

| 状态          | 说明       |
| ----------- | -------- |
| token       | 登录 token |
| currentUser | 当前登录用户   |
| isLogin     | 是否登录     |

---

## 20.2 appStore

```ts
interface AppState {
  unreadNotificationCount: number;
  setUnreadNotificationCount: (count: number) => void;
}
```

职责：

| 状态                      | 说明       |
| ----------------------- | -------- |
| unreadNotificationCount | 顶部通知红点数量 |

---

# 21. API 文件拆分建议

## 21.1 auth.ts

```ts
export function login(data: LoginRequest) {}
export function register(data: RegisterRequest) {}
export function getCurrentUser() {}
export function logout() {}
```

## 21.2 post.ts

```ts
export function getPostList(params: PostQuery) {}
export function getFollowingPostList(params: PageRequest) {}
export function getHotPosts(params: { type: 'today' | 'week' | 'all'; limit: number }) {}
export function getPostDetail(postId: number) {}
export function createPost(data: CreatePostRequest) {}
export function updatePost(postId: number, data: UpdatePostRequest) {}
export function deletePost(postId: number) {}
```

## 21.3 comment.ts

```ts
export function getPostComments(postId: number, params: PageRequest) {}
export function createComment(data: CreateCommentRequest) {}
export function createReply(commentId: number, data: CreateReplyRequest) {}
export function deleteComment(commentId: number) {}
```

## 21.4 interaction.ts

```ts
export function likePost(postId: number) {}
export function unlikePost(postId: number) {}
export function favoritePost(postId: number) {}
export function unfavoritePost(postId: number) {}
export function likeComment(commentId: number) {}
export function unlikeComment(commentId: number) {}
```

## 21.5 follow.ts

```ts
export function followUser(userId: number) {}
export function unfollowUser(userId: number) {}
export function blockUser(userId: number) {}
```

## 21.6 notification.ts

```ts
export function getNotifications(params: {
  type: string;
  page: number;
  pageSize: number;
}) {}

export function readNotification(notificationId: number) {}
export function readAllNotifications() {}
export function getUnreadNotificationCount() {}
```

## 21.7 report.ts

```ts
export function createReport(data: CreateReportRequest) {}
```

## 21.8 admin.ts

```ts
export function getAdminDashboardStats() {}
export function getAdminUsers(params: any) {}
export function muteUser(userId: number) {}
export function banUser(userId: number) {}
export function unbanUser(userId: number) {}

export function getAdminPosts(params: any) {}
export function approvePost(postId: number) {}
export function rejectPost(postId: number, reason: string) {}
export function takedownPost(postId: number, reason: string) {}

export function getAdminComments(params: any) {}
export function deleteAdminComment(commentId: number) {}

export function getAdminReports(params: any) {}
export function approveReport(reportId: number, data: any) {}
export function rejectReport(reportId: number, data: any) {}

export function getSensitiveWords(params: any) {}
export function createSensitiveWord(data: any) {}
export function updateSensitiveWord(wordId: number, data: any) {}
export function deleteSensitiveWord(wordId: number) {}

export function getOperationLogs(params: any) {}
```

---

# 22. 页面与接口对应关系

## 22.1 前台页面

| 页面    | 依赖接口                                                                |
| ----- | ------------------------------------------------------------------- |
| 登录页   | `POST /auth/login`                                                  |
| 注册页   | `POST /auth/register`                                               |
| 首页    | `GET /posts`、点赞、收藏、关注                                               |
| 热门页   | `GET /posts/hot`                                                    |
| 关注页   | `GET /posts/following`                                              |
| 帖子详情页 | `GET /posts/{postId}`、`GET /posts/{postId}/comments`、评论、回复、点赞、收藏、举报 |
| 发布帖子页 | `POST /posts`、`POST /upload/image`                                  |
| 编辑帖子页 | `GET /posts/{postId}`、`PUT /posts/{postId}`、`POST /upload/image`    |
| 我的主页  | `GET /auth/me`、我的帖子、我的点赞、我的收藏、我的关注                                  |
| 用户主页  | `GET /users/{userId}`、`GET /users/{userId}/posts`、关注、拉黑、举报          |
| 通知页   | `GET /notifications`、已读接口                                           |

---

## 22.2 后台页面

| 页面    | 依赖接口                                                           |
| ----- | -------------------------------------------------------------- |
| 数据看板  | `GET /admin/dashboard/stats`                                   |
| 用户管理  | `GET /admin/users`、禁言、封禁、解封                                    |
| 帖子管理  | `GET /admin/posts`、审核、下架、删除                                    |
| 评论管理  | `GET /admin/comments`、删除评论                                     |
| 审核管理  | `GET /admin/reviews/posts`、`GET /admin/reviews/comments`、通过、拒绝 |
| 举报管理  | `GET /admin/reports`、通过举报、驳回举报                                 |
| 敏感词管理 | 敏感词增删改查                                                        |
| 操作日志  | `GET /admin/operation-logs`                                    |

---

# 23. 前端 Mock 实现建议

开发初期可以在 `api/mock.ts` 中实现假接口。

## 23.1 模拟延迟

```ts
export function mockDelay<T>(data: T, delay = 300): Promise<T> {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(data);
    }, delay);
  });
}
```

## 23.2 模拟分页

```ts
export function mockPage<T>(
  list: T[],
  page: number,
  pageSize: number,
): PageResult<T> {
  const start = (page - 1) * pageSize;
  const end = start + pageSize;

  return {
    list: list.slice(start, end),
    total: list.length,
    page,
    pageSize,
  };
}
```

## 23.3 模拟点赞

```ts
export function togglePostLike(postId: number) {
  const post = mockPosts.find((item) => item.postId === postId);

  if (!post) {
    throw new Error('帖子不存在');
  }

  post.liked = !post.liked;
  post.likeCount += post.liked ? 1 : -1;

  return {
    liked: post.liked,
    likeCount: post.likeCount,
  };
}
```

---

# 24. 组件实现优先级

## P0：优先实现

```text
FrontLayout
LoginPage
RegisterPage
HomePage
PostCard
PostList
PostDetailPage
CommentInput
CommentList
PostCreatePage
MyProfilePage
AdminLayout
AdminDashboardPage
AdminUserPage
AdminPostPage
```

## P1：继续实现

```text
HotPage
FollowingPage
PostEditPage
UserProfilePage
NotificationPage
ReportModal
AdminCommentPage
AdminReviewPage
AdminReportPage
```

## P2：最后完善

```text
AdminSensitiveWordPage
AdminOperationLogPage
ImagePreviewModal
UploadImageGrid
各类 Drawer / Modal 细节
错误状态
骨架屏
```

---

# 25. 给智能编程助手的最终实现提示词

请根据以下要求实现社区 MVP 前端项目：

```text
使用 React + TypeScript + React Router + Zustand + Axios + Ant Design 实现一个社区 MVP 前端项目。

需要实现前台和后台两套布局。

前台页面包括：
1. 登录页
2. 注册页
3. 首页帖子流
4. 热门榜单页
5. 关注流页面
6. 帖子详情页
7. 发布帖子页
8. 编辑帖子页
9. 我的主页
10. 其他用户主页
11. 通知页

后台页面包括：
1. 数据看板
2. 用户管理
3. 帖子管理
4. 评论管理
5. 审核管理
6. 举报管理
7. 敏感词管理
8. 操作日志

项目需要：
1. 按照文档中的组件拆分进行开发。
2. 使用 mock 数据模拟后端接口。
3. 所有列表需要支持 loading、空状态、错误状态、分页。
4. 表单需要支持基础校验。
5. 点赞、收藏、关注需要能即时更新 UI 状态。
6. 后台表格操作需要弹窗确认。
7. 页面风格需要统一：浅灰背景、白色卡片、圆角、轻阴影、蓝色主色。
8. API 文件需要按照 auth、post、comment、interaction、follow、notification、report、admin 拆分。
9. 类型定义需要放在 types 目录中。
10. 路由需要按照文档中的路由约定实现。
```
