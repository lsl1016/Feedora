# 增强版社区 MVP 前端组件拆分清单 + Mock 数据结构 + API 接口约定

## 1. 技术栈建议

推荐技术栈：

```text
React + TypeScript + React Router + Zustand + Axios + Ant Design + Vite
```

推荐原因：

```text
1. React + TypeScript 适合组件化开发。
2. React Router 负责前后台路由。
3. Zustand 负责轻量状态管理。
4. Axios 统一封装请求。
5. Ant Design 提供表单、表格、弹窗、上传、菜单、抽屉等后台和前台常用组件。
6. Vite 启动快，适合快速开发原型。
```

---

# 2. 项目目录结构建议

```text
community-mvp-enhanced-frontend
├── public
│   └── logo.svg
│
├── src
│   ├── api
│   │   ├── request.ts
│   │   ├── mockAdapter.ts
│   │   ├── auth.ts
│   │   ├── user.ts
│   │   ├── post.ts
│   │   ├── comment.ts
│   │   ├── interaction.ts
│   │   ├── follow.ts
│   │   ├── notification.ts
│   │   ├── report.ts
│   │   ├── tag.ts
│   │   ├── circle.ts
│   │   ├── topic.ts
│   │   ├── search.ts
│   │   ├── recommend.ts
│   │   ├── growth.ts
│   │   ├── activity.ts
│   │   ├── announcement.ts
│   │   ├── upload.ts
│   │   └── admin.ts
│   │
│   ├── mock
│   │   ├── index.ts
│   │   ├── user.mock.ts
│   │   ├── post.mock.ts
│   │   ├── comment.mock.ts
│   │   ├── notification.mock.ts
│   │   ├── tag.mock.ts
│   │   ├── circle.mock.ts
│   │   ├── topic.mock.ts
│   │   ├── search.mock.ts
│   │   ├── recommend.mock.ts
│   │   ├── growth.mock.ts
│   │   ├── activity.mock.ts
│   │   ├── announcement.mock.ts
│   │   └── admin.mock.ts
│   │
│   ├── assets
│   │   ├── images
│   │   └── icons
│   │
│   ├── components
│   │   ├── common
│   │   ├── layout
│   │   ├── post
│   │   ├── comment
│   │   ├── user
│   │   ├── tag
│   │   ├── circle
│   │   ├── topic
│   │   ├── search
│   │   ├── notification
│   │   ├── growth
│   │   ├── activity
│   │   ├── report
│   │   └── admin
│   │
│   ├── constants
│   │   ├── enums.ts
│   │   ├── route.ts
│   │   ├── options.ts
│   │   └── permissions.ts
│   │
│   ├── hooks
│   │   ├── useAuth.ts
│   │   ├── usePagination.ts
│   │   ├── usePostActions.ts
│   │   ├── useFollowActions.ts
│   │   ├── useCircleActions.ts
│   │   ├── useSearchSuggest.ts
│   │   └── useAdminTable.ts
│   │
│   ├── layouts
│   │   ├── FrontLayout.tsx
│   │   └── AdminLayout.tsx
│   │
│   ├── pages
│   │   ├── auth
│   │   ├── home
│   │   ├── hot
│   │   ├── post
│   │   ├── user
│   │   ├── tag
│   │   ├── circle
│   │   ├── topic
│   │   ├── search
│   │   ├── notification
│   │   ├── growth
│   │   ├── activity
│   │   ├── announcement
│   │   └── admin
│   │
│   ├── router
│   │   └── index.tsx
│   │
│   ├── store
│   │   ├── authStore.ts
│   │   ├── appStore.ts
│   │   ├── userStore.ts
│   │   └── draftStore.ts
│   │
│   ├── styles
│   │   ├── global.css
│   │   ├── variables.css
│   │   └── admin.css
│   │
│   ├── types
│   │   ├── common.ts
│   │   ├── user.ts
│   │   ├── post.ts
│   │   ├── comment.ts
│   │   ├── interaction.ts
│   │   ├── notification.ts
│   │   ├── report.ts
│   │   ├── tag.ts
│   │   ├── circle.ts
│   │   ├── topic.ts
│   │   ├── search.ts
│   │   ├── recommend.ts
│   │   ├── growth.ts
│   │   ├── activity.ts
│   │   ├── announcement.ts
│   │   └── admin.ts
│   │
│   ├── utils
│   │   ├── format.ts
│   │   ├── storage.ts
│   │   ├── permission.ts
│   │   ├── validator.ts
│   │   ├── mock.ts
│   │   └── clipboard.ts
│   │
│   ├── App.tsx
│   └── main.tsx
│
├── package.json
├── vite.config.ts
├── tsconfig.json
└── README.md
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
├── post
│   ├── PostDetailPage.tsx
│   ├── PostCreatePage.tsx
│   └── PostEditPage.tsx
│
├── user
│   ├── MyProfilePage.tsx
│   ├── UserProfilePage.tsx
│   └── CreatorCenterPage.tsx
│
├── tag
│   └── TagDetailPage.tsx
│
├── circle
│   ├── CircleListPage.tsx
│   ├── CircleCreatePage.tsx
│   ├── CircleHomePage.tsx
│   └── CircleMemberManagePage.tsx
│
├── topic
│   ├── TopicListPage.tsx
│   └── TopicDetailPage.tsx
│
├── search
│   └── SearchResultPage.tsx
│
├── notification
│   └── NotificationPage.tsx
│
├── growth
│   ├── TaskCenterPage.tsx
│   └── RankingPage.tsx
│
├── activity
│   ├── ActivityCenterPage.tsx
│   └── ActivityDetailPage.tsx
│
└── announcement
    ├── AnnouncementListPage.tsx
    └── AnnouncementDetailPage.tsx
```

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
    ├── AdminTagPage.tsx
    ├── AdminCirclePage.tsx
    ├── AdminTopicPage.tsx
    ├── AdminSearchPage.tsx
    ├── AdminRecommendPage.tsx
    ├── AdminPointLevelPage.tsx
    ├── AdminBadgePage.tsx
    ├── AdminTaskPage.tsx
    ├── AdminActivityPage.tsx
    ├── AdminBannerPage.tsx
    ├── AdminAnnouncementPage.tsx
    ├── AdminSensitiveWordPage.tsx
    ├── AdminOperationLogPage.tsx
    └── AdminOperationDashboardPage.tsx
```

---

# 4. 布局组件拆分

## 4.1 前台布局

```text
layouts
└── FrontLayout.tsx
```

组件结构：

```text
FrontLayout
├── FrontHeader
│   ├── Logo
│   ├── FrontNavMenu
│   ├── GlobalSearchBox
│   ├── PublishButton
│   ├── NotificationEntry
│   └── UserAvatarDropdown
└── FrontMainContent
```

组件职责：

| 组件                   | 说明                    |
| -------------------- | --------------------- |
| `FrontHeader`        | 前台顶部导航栏               |
| `Logo`               | 社区名称，点击跳转首页           |
| `FrontNavMenu`       | 首页、推荐、热门、圈子、话题、活动、排行榜 |
| `GlobalSearchBox`    | 全局搜索框，支持热门搜索和搜索联想     |
| `PublishButton`      | 发布帖子入口                |
| `NotificationEntry`  | 通知入口，展示未读数            |
| `UserAvatarDropdown` | 用户头像下拉菜单              |
| `FrontMainContent`   | 页面内容容器                |

---

## 4.2 后台布局

```text
layouts
└── AdminLayout.tsx
```

组件结构：

```text
AdminLayout
├── AdminSidebar
├── AdminTopbar
├── AdminBreadcrumb
└── AdminContent
```

组件职责：

| 组件                | 说明     |
| ----------------- | ------ |
| `AdminSidebar`    | 后台左侧菜单 |
| `AdminTopbar`     | 后台顶部栏  |
| `AdminBreadcrumb` | 面包屑导航  |
| `AdminContent`    | 后台内容区域 |

---

# 5. 通用组件拆分

```text
components
└── common
    ├── PageContainer.tsx
    ├── CardContainer.tsx
    ├── SectionTitle.tsx
    ├── EmptyState.tsx
    ├── ErrorState.tsx
    ├── LoadingSkeleton.tsx
    ├── ConfirmModal.tsx
    ├── StatusTag.tsx
    ├── ImagePreviewModal.tsx
    ├── UploadImageGrid.tsx
    ├── FormErrorText.tsx
    ├── PaginationBar.tsx
    ├── CopyButton.tsx
    ├── RichTextPreview.tsx
    ├── DataMetricCard.tsx
    ├── FilterBar.tsx
    └── ActionDropdown.tsx
```

组件说明：

| 组件                  | 说明             |
| ------------------- | -------------- |
| `PageContainer`     | 页面统一容器，控制宽度和背景 |
| `CardContainer`     | 白色卡片容器         |
| `SectionTitle`      | 区块标题           |
| `EmptyState`        | 空状态            |
| `ErrorState`        | 错误状态           |
| `LoadingSkeleton`   | 骨架屏            |
| `ConfirmModal`      | 通用确认弹窗         |
| `StatusTag`         | 状态标签           |
| `ImagePreviewModal` | 图片预览弹窗         |
| `UploadImageGrid`   | 图片上传宫格         |
| `FormErrorText`     | 表单错误提示         |
| `PaginationBar`     | 分页组件           |
| `CopyButton`        | 复制按钮           |
| `RichTextPreview`   | 正文预览           |
| `DataMetricCard`    | 指标卡片           |
| `FilterBar`         | 通用筛选区          |
| `ActionDropdown`    | 更多操作下拉菜单       |

---

# 6. 前台业务组件拆分

## 6.1 帖子组件

```text
components
└── post
    ├── PostCard.tsx
    ├── RepostCard.tsx
    ├── SourcePostCard.tsx
    ├── PostList.tsx
    ├── PostFeedTabs.tsx
    ├── PostFilterBar.tsx
    ├── PublishEntryCard.tsx
    ├── PostActionBar.tsx
    ├── PostAuthorInfo.tsx
    ├── PostContent.tsx
    ├── PostForm.tsx
    ├── PostTagSelector.tsx
    ├── PostTopicSelector.tsx
    ├── PostCircleSelector.tsx
    ├── SchedulePublishSelector.tsx
    ├── PostPreviewModal.tsx
    ├── ShareModal.tsx
    ├── RepostModal.tsx
    ├── SimilarPostList.tsx
    ├── HotRankList.tsx
    └── HotRankItem.tsx
```

说明：

| 组件                        | 说明                 |
| ------------------------- | ------------------ |
| `PostCard`                | 普通帖子卡片             |
| `RepostCard`              | 转发帖子卡片             |
| `SourcePostCard`          | 转发中的原帖卡片           |
| `PostList`                | 帖子列表               |
| `PostFeedTabs`            | 推荐 / 最新 / 热门 / 关注  |
| `PostFilterBar`           | 发布时间、标签、圈子、话题、排序筛选 |
| `PublishEntryCard`        | 首页发布入口             |
| `PostActionBar`           | 点赞、评论、收藏、分享、转发、举报  |
| `PostAuthorInfo`          | 作者信息               |
| `PostContent`             | 帖子正文               |
| `PostForm`                | 发布 / 编辑帖子表单        |
| `PostTagSelector`         | 内容标签选择器            |
| `PostTopicSelector`       | 话题选择器              |
| `PostCircleSelector`      | 圈子选择器              |
| `SchedulePublishSelector` | 定时发布选择器            |
| `PostPreviewModal`        | 帖子预览弹窗             |
| `ShareModal`              | 分享弹窗               |
| `RepostModal`             | 转发弹窗               |
| `SimilarPostList`         | 相关推荐               |
| `HotRankList`             | 热门榜单               |
| `HotRankItem`             | 榜单项                |

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

---

## 6.3 用户组件

```text
components
└── user
    ├── UserProfileCard.tsx
    ├── UserStatBar.tsx
    ├── UserLevelBadge.tsx
    ├── EditProfileModal.tsx
    ├── FollowButton.tsx
    ├── UserList.tsx
    ├── UserListItem.tsx
    ├── CheckInButton.tsx
    ├── CheckInModal.tsx
    ├── BadgeWall.tsx
    └── CreatorStatsPanel.tsx
```

---

## 6.4 标签组件

```text
components
└── tag
    ├── TagChip.tsx
    ├── TagList.tsx
    ├── TagSelector.tsx
    └── TagHeader.tsx
```

---

## 6.5 圈子组件

```text
components
└── circle
    ├── CircleCard.tsx
    ├── CircleList.tsx
    ├── CircleFilterBar.tsx
    ├── CircleHeader.tsx
    ├── CircleTabs.tsx
    ├── CirclePostList.tsx
    ├── CircleAnnouncementList.tsx
    ├── CircleAnnouncementModal.tsx
    ├── CircleRulePanel.tsx
    ├── CircleMemberList.tsx
    ├── CircleMemberTable.tsx
    ├── CircleJoinButton.tsx
    ├── CircleJoinRequestModal.tsx
    ├── CircleMemberMuteModal.tsx
    └── CircleRoleActionMenu.tsx
```

---

## 6.6 话题组件

```text
components
└── topic
    ├── TopicCard.tsx
    ├── TopicList.tsx
    ├── TopicFilterTabs.tsx
    ├── TopicHeader.tsx
    └── TopicPostList.tsx
```

---

## 6.7 搜索组件

```text
components
└── search
    ├── GlobalSearchBox.tsx
    ├── SearchSuggestDropdown.tsx
    ├── HotSearchList.tsx
    ├── SearchResultTabs.tsx
    ├── SearchResultList.tsx
    ├── SearchPostResultItem.tsx
    ├── SearchCommentResultItem.tsx
    ├── SearchUserResultItem.tsx
    ├── SearchTopicResultItem.tsx
    └── SearchCircleResultItem.tsx
```

---

## 6.8 通知组件

```text
components
└── notification
    ├── NotificationTabs.tsx
    ├── NotificationList.tsx
    ├── NotificationItem.tsx
    └── NotificationDetailModal.tsx
```

---

## 6.9 成长体系组件

```text
components
└── growth
    ├── PointSummaryCard.tsx
    ├── LevelProgressCard.tsx
    ├── TaskTabs.tsx
    ├── TaskCard.tsx
    ├── RankingTabs.tsx
    ├── RankingTable.tsx
    ├── BadgeCard.tsx
    └── BadgeList.tsx
```

---

## 6.10 活动组件

```text
components
└── activity
    ├── ActivityBanner.tsx
    ├── ActivityTabs.tsx
    ├── ActivityCard.tsx
    ├── ActivityList.tsx
    ├── TopicActivityDetail.tsx
    ├── VoteActivityDetail.tsx
    ├── EssayActivityDetail.tsx
    ├── CheckInActivityDetail.tsx
    ├── VoteOptionList.tsx
    └── ActivitySubmissionList.tsx
```

---

## 6.11 举报组件

```text
components
└── report
    └── ReportModal.tsx
```

---

# 7. 后台业务组件拆分

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
    ├── TagFormModal.tsx
    ├── CircleDetailDrawer.tsx
    ├── CircleCloseModal.tsx
    ├── TopicFormModal.tsx
    ├── RecommendContentModal.tsx
    ├── PointRuleFormModal.tsx
    ├── LevelRuleFormModal.tsx
    ├── BadgeFormModal.tsx
    ├── BadgeGrantModal.tsx
    ├── TaskFormModal.tsx
    ├── ActivityFormModal.tsx
    ├── ActivityDataDrawer.tsx
    ├── BannerFormModal.tsx
    ├── AnnouncementFormModal.tsx
    ├── SensitiveWordFormModal.tsx
    └── OperationLogDetailDrawer.tsx
```

---

# 8. TypeScript 通用类型定义

## 8.1 common.ts

```ts
export type ID = number;

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

export type TargetType =
  | 'post'
  | 'comment'
  | 'user'
  | 'tag'
  | 'circle'
  | 'topic'
  | 'activity'
  | 'announcement';

export type SortType =
  | 'recommend'
  | 'latest'
  | 'hot'
  | 'comment'
  | 'favorite'
  | 'view';

export type EnabledStatus = 'enabled' | 'disabled';

export type TimeRange = 'all' | 'today' | 'week' | 'month';

export interface SelectOption {
  label: string;
  value: string | number;
}
```

---

# 9. 用户类型定义

## 9.1 user.ts

```ts
export type UserStatus = 'normal' | 'muted' | 'banned' | 'deleted';

export type UserRole = 'user' | 'admin';

export interface User {
  userId: number;
  account: string;
  nickname: string;
  avatar: string;
  bio: string;
  role: UserRole;
  status: UserStatus;

  level: number;
  levelName: string;
  points: number;
  experience: number;
  nextLevelExperience: number;

  postCount: number;
  commentCount: number;
  followerCount: number;
  followingCount: number;
  likeReceivedCount: number;
  badgeCount: number;

  checkedInToday: boolean;
  continuousCheckInDays: number;

  createdAt: string;
  updatedAt: string;
}

export interface UserSummary {
  userId: number;
  nickname: string;
  avatar: string;
  bio?: string;
  level?: number;
  levelName?: string;
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

# 10. 标签类型定义

## 10.1 tag.ts

```ts
import type { EnabledStatus } from './common';

export interface ContentTag {
  tagId: number;
  tagName: string;
  description: string;
  status: EnabledStatus;
  useCount: number;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface TagSummary {
  tagId: number;
  tagName: string;
}

export interface TagQuery {
  keyword?: string;
  status?: EnabledStatus | 'all';
  page?: number;
  pageSize?: number;
}

export interface CreateTagRequest {
  tagName: string;
  description: string;
  status: EnabledStatus;
}

export interface UpdateTagRequest {
  tagName: string;
  description: string;
  status: EnabledStatus;
}
```

---

# 11. 圈子类型定义

## 11.1 circle.ts

```ts
import type { EnabledStatus } from './common';
import type { TagSummary } from './tag';
import type { UserSummary } from './user';

export type CircleStatus = 'normal' | 'reviewing' | 'closed' | 'deleted';

export type CircleJoinType = 'direct' | 'approval';

export type CirclePostPermission = 'all' | 'admin_only';

export type CircleMemberRole = 'owner' | 'moderator' | 'reviewer' | 'member';

export type CircleMemberStatus = 'normal' | 'muted' | 'removed';

export type CircleJoinRequestStatus = 'pending' | 'approved' | 'rejected';

export interface Circle {
  circleId: number;
  name: string;
  avatar: string;
  description: string;
  category: string;
  tags: TagSummary[];
  ownerId: number;
  owner: UserSummary;
  joinType: CircleJoinType;
  postPermission: CirclePostPermission;
  memberCount: number;
  postCount: number;
  featuredPostCount: number;
  isJoined: boolean;
  myRole?: CircleMemberRole;
  myStatus?: CircleMemberStatus;
  isRecommended: boolean;
  status: CircleStatus;
  rules: string;
  createdAt: string;
  updatedAt: string;
}

export interface CircleMember {
  id: number;
  circleId: number;
  userId: number;
  user: UserSummary;
  role: CircleMemberRole;
  status: CircleMemberStatus;
  muteReason?: string;
  mutedUntil?: string;
  joinedAt: string;
  updatedAt: string;
}

export interface CircleJoinRequest {
  requestId: number;
  circleId: number;
  userId: number;
  user: UserSummary;
  reason: string;
  status: CircleJoinRequestStatus;
  handledBy?: number;
  handleReason?: string;
  createdAt: string;
  handledAt?: string;
}

export interface CircleAnnouncement {
  announcementId: number;
  circleId: number;
  title: string;
  content: string;
  publisherId: number;
  publisher: UserSummary;
  status: EnabledStatus;
  createdAt: string;
  updatedAt: string;
}

export interface CreateCircleRequest {
  name: string;
  avatar: string;
  description: string;
  category: string;
  tagIds: number[];
  joinType: CircleJoinType;
  postPermission: CirclePostPermission;
  rules: string;
}

export interface CircleQuery {
  keyword?: string;
  category?: string;
  tagId?: number;
  sort?: 'recommend' | 'latest' | 'member' | 'post';
  status?: CircleStatus | 'all';
  page: number;
  pageSize: number;
}
```

---

# 12. 话题类型定义

## 12.1 topic.ts

```ts
import type { EnabledStatus } from './common';

export interface Topic {
  topicId: number;
  name: string;
  description: string;
  coverImage: string;
  postCount: number;
  participantCount: number;
  isOfficial: boolean;
  isRecommended: boolean;
  status: EnabledStatus;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface TopicSummary {
  topicId: number;
  name: string;
}

export interface TopicQuery {
  keyword?: string;
  type?: 'all' | 'official' | 'hot' | 'latest';
  status?: EnabledStatus | 'all';
  page: number;
  pageSize: number;
}

export interface CreateTopicRequest {
  name: string;
  description: string;
  coverImage: string;
  isOfficial: boolean;
  isRecommended: boolean;
  status: EnabledStatus;
}
```

---

# 13. 帖子类型定义

## 13.1 post.ts

```ts
import type { SortType, TimeRange } from './common';
import type { UserSummary } from './user';
import type { TagSummary } from './tag';
import type { TopicSummary } from './topic';
import type { Circle } from './circle';

export type PostStatus =
  | 'draft'
  | 'scheduled'
  | 'reviewing'
  | 'published'
  | 'rejected'
  | 'deleted'
  | 'takedown';

export type PostType = 'original' | 'repost';

export interface Post {
  postId: number;
  postType: PostType;

  authorId: number;
  author: UserSummary;

  title: string;
  content: string;
  summary: string;
  images: string[];

  tags: TagSummary[];
  topics: TopicSummary[];

  circleId?: number;
  circle?: Pick<Circle, 'circleId' | 'name' | 'avatar' | 'memberCount' | 'postCount'>;

  status: PostStatus;
  rejectReason?: string;
  takedownReason?: string;

  isTop: boolean;
  isFeatured: boolean;
  isSelected: boolean;

  scheduledAt?: string;

  sourcePostId?: number;
  sourcePost?: Post;
  repostComment?: string;

  viewCount: number;
  likeCount: number;
  commentCount: number;
  favoriteCount: number;
  shareCount: number;
  repostCount: number;
  hotScore: number;

  liked: boolean;
  favorited: boolean;
  followedAuthor: boolean;

  createdAt: string;
  publishedAt?: string;
  updatedAt: string;
}

export interface CreatePostRequest {
  title: string;
  content: string;
  images: string[];
  tagIds: number[];
  topicIds: number[];
  circleId?: number;
  publishMode: 'now' | 'schedule' | 'draft';
  scheduledAt?: string;
}

export interface UpdatePostRequest {
  title: string;
  content: string;
  images: string[];
  tagIds: number[];
  topicIds: number[];
  circleId?: number;
  publishMode?: 'now' | 'schedule' | 'draft';
  scheduledAt?: string;
}

export interface RepostRequest {
  sourcePostId: number;
  repostComment: string;
}

export interface PostQuery {
  page: number;
  pageSize: number;
  feedType?: 'recommend' | 'latest' | 'hot' | 'following';
  sort?: SortType;
  timeRange?: TimeRange;
  tagId?: number;
  circleId?: number;
  topicId?: number;
  keyword?: string;
  status?: PostStatus;
  authorId?: number;
}
```

---

# 14. 评论类型定义

## 14.1 comment.ts

```ts
import type { UserSummary } from './user';

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

# 15. 通知类型定义

## 15.1 notification.ts

```ts
import type { TargetType } from './common';
import type { UserSummary } from './user';

export type NotificationType =
  | 'like_post'
  | 'like_comment'
  | 'comment_post'
  | 'reply_comment'
  | 'follow_user'
  | 'review_pass'
  | 'review_reject'
  | 'system'
  | 'circle_join_request'
  | 'circle_join_approved'
  | 'circle_join_rejected'
  | 'circle_muted'
  | 'circle_removed'
  | 'circle_role_changed'
  | 'post_featured'
  | 'post_top'
  | 'scheduled_post_success'
  | 'scheduled_post_failed'
  | 'point_changed'
  | 'badge_granted'
  | 'activity_reminder'
  | 'official_announcement';

export type NotificationCategory =
  | 'all'
  | 'interaction'
  | 'follow'
  | 'circle'
  | 'review'
  | 'growth'
  | 'activity'
  | 'system';

export type ReadStatus = 'unread' | 'read';

export interface Notification {
  notificationId: number;
  receiverId: number;
  senderId?: number;
  sender?: UserSummary;
  type: NotificationType;
  category: NotificationCategory;
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

# 16. 搜索与推荐类型定义

## 16.1 search.ts

```ts
import type { TargetType } from './common';
import type { Post } from './post';
import type { User } from './user';
import type { Topic } from './topic';
import type { Circle } from './circle';

export type SearchResultType = 'all' | 'post' | 'comment' | 'user' | 'topic' | 'circle';

export interface SearchSuggestItem {
  type: TargetType;
  id: number;
  title: string;
  description?: string;
}

export interface HotSearchKeyword {
  keywordId: number;
  keyword: string;
  searchCount: number;
  status: 'normal' | 'blocked';
  lastSearchAt: string;
}

export interface SearchResult {
  posts: Post[];
  users: User[];
  topics: Topic[];
  circles: Circle[];
  comments: SearchCommentResult[];
}

export interface SearchCommentResult {
  commentId: number;
  postId: number;
  postTitle: string;
  content: string;
  userId: number;
  nickname: string;
  createdAt: string;
}
```

## 16.2 recommend.ts

```ts
import type { Post } from './post';
import type { Circle } from './circle';
import type { Topic } from './topic';

export interface RecommendContentPoolItem {
  id: number;
  targetType: 'post' | 'circle' | 'topic' | 'activity';
  targetId: number;
  title: string;
  weight: number;
  status: 'online' | 'offline';
  createdAt: string;
}

export interface RecommendResult {
  posts: Post[];
  circles: Circle[];
  topics: Topic[];
}
```

---

# 17. 成长体系类型定义

## 17.1 growth.ts

```ts
export type TaskType = 'newbie' | 'daily' | 'growth';

export type TaskStatus = 'todo' | 'done' | 'claimed';

export interface PointLog {
  logId: number;
  userId: number;
  points: number;
  reason: string;
  createdAt: string;
}

export interface LevelRule {
  level: number;
  levelName: string;
  minExperience: number;
  maxExperience: number;
  benefits: string[];
}

export interface Badge {
  badgeId: number;
  name: string;
  icon: string;
  description: string;
  condition: string;
  status: 'enabled' | 'disabled';
  createdAt: string;
}

export interface UserBadge {
  id: number;
  userId: number;
  badge: Badge;
  grantedAt: string;
}

export interface Task {
  taskId: number;
  title: string;
  description: string;
  type: TaskType;
  rewardPoints: number;
  rewardBadgeId?: number;
  targetValue: number;
  currentValue: number;
  status: TaskStatus;
  actionText: string;
  actionUrl: string;
}

export interface CheckInResult {
  checkedIn: boolean;
  points: number;
  continuousDays: number;
}

export type RankingType = 'active' | 'contribution' | 'creator' | 'circle';

export type RankingTimeRange = 'day' | 'week' | 'all';

export interface RankingItem {
  rank: number;
  targetId: number;
  targetType: 'user' | 'circle';
  name: string;
  avatar: string;
  level?: number;
  levelName?: string;
  points?: number;
  postCount?: number;
  likeReceivedCount?: number;
  memberCount?: number;
  score: number;
  isCurrentUser?: boolean;
}
```

---

# 18. 活动与公告类型定义

## 18.1 activity.ts

```ts
export type ActivityType = 'topic' | 'vote' | 'essay' | 'checkin';

export type ActivityStatus = 'not_started' | 'ongoing' | 'ended' | 'offline';

export interface Activity {
  activityId: number;
  type: ActivityType;
  title: string;
  description: string;
  coverImage: string;
  topicId?: number;
  topicName?: string;
  startAt: string;
  endAt: string;
  participantCount: number;
  submissionCount: number;
  status: ActivityStatus;

  voteConfig?: VoteConfig;
  checkInConfig?: CheckInConfig;
}

export interface VoteConfig {
  question: string;
  multiple: boolean;
  options: VoteOption[];
  voted: boolean;
  selectedOptionIds: number[];
}

export interface VoteOption {
  optionId: number;
  text: string;
  voteCount: number;
  percent: number;
}

export interface CheckInConfig {
  totalDays: number;
  currentDays: number;
  todayChecked: boolean;
  todayTask: string;
}

export interface ActivitySubmission {
  submissionId: number;
  activityId: number;
  postId: number;
  postTitle: string;
  userId: number;
  nickname: string;
  createdAt: string;
}
```

## 18.2 announcement.ts

```ts
export interface OfficialAnnouncement {
  announcementId: number;
  title: string;
  summary: string;
  content: string;
  publisherId: number;
  publisherName: string;
  status: 'draft' | 'published' | 'offline';
  publishedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Banner {
  bannerId: number;
  title: string;
  imageUrl: string;
  targetUrl: string;
  position: 'home' | 'hot' | 'circle' | 'topic' | 'activity';
  weight: number;
  status: 'enabled' | 'disabled';
  createdAt: string;
}
```

---

# 19. 举报类型定义

## 19.1 report.ts

```ts
import type { TargetType } from './common';
import type { UserSummary } from './user';

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

# 20. 后台类型定义

## 20.1 admin.ts

```ts
export interface AdminDashboardStats {
  todayNewUserCount: number;
  todayPostCount: number;
  todayCommentCount: number;
  todayLikeCount: number;
  pendingReviewCount: number;
  pendingReportCount: number;
  pendingCircleCount: number;
  todayActivityParticipantCount: number;
}

export interface OperationDashboardStats {
  dau: number;
  newUserCount: number;
  postCount: number;
  commentCount: number;
  activityParticipantCount: number;
  retentionRate: number;
  conversionRate: number;
}

export interface TrendPoint {
  date: string;
  value: number;
}

export type SensitiveWordStrategy = 'block' | 'review';

export interface SensitiveWord {
  wordId: number;
  word: string;
  strategy: SensitiveWordStrategy;
  status: 'enabled' | 'disabled';
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

# 21. Mock 数据结构

## 21.1 mockUsers

```ts
export const mockUsers: User[] = [
  {
    userId: 10001,
    account: 'zhangsan',
    nickname: '张三',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan',
    bio: '后端开发，喜欢 Go、Java 和数据平台。',
    role: 'user',
    status: 'normal',
    level: 3,
    levelName: '内容贡献者',
    points: 560,
    experience: 720,
    nextLevelExperience: 800,
    postCount: 12,
    commentCount: 35,
    followerCount: 256,
    followingCount: 86,
    likeReceivedCount: 1234,
    badgeCount: 3,
    checkedInToday: false,
    continuousCheckInDays: 2,
    createdAt: '2026-07-01 10:00:00',
    updatedAt: '2026-07-05 10:00:00',
  },
  {
    userId: 10002,
    account: 'frontend_dev',
    nickname: '前端小李',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=frontend',
    bio: 'React / Vue 前端开发者，喜欢做 UI 原型。',
    role: 'user',
    status: 'normal',
    level: 4,
    levelName: '社区达人',
    points: 980,
    experience: 1200,
    nextLevelExperience: 2000,
    postCount: 18,
    commentCount: 66,
    followerCount: 512,
    followingCount: 120,
    likeReceivedCount: 3450,
    badgeCount: 5,
    checkedInToday: true,
    continuousCheckInDays: 7,
    createdAt: '2026-07-02 11:20:00',
    updatedAt: '2026-07-05 12:00:00',
  },
  {
    userId: 10003,
    account: 'admin',
    nickname: '社区管理员',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin',
    bio: '社区运营与内容治理。',
    role: 'admin',
    status: 'normal',
    level: 5,
    levelName: '核心创作者',
    points: 2400,
    experience: 2600,
    nextLevelExperience: 3000,
    postCount: 6,
    commentCount: 12,
    followerCount: 900,
    followingCount: 30,
    likeReceivedCount: 5200,
    badgeCount: 8,
    checkedInToday: true,
    continuousCheckInDays: 15,
    createdAt: '2026-06-30 09:00:00',
    updatedAt: '2026-07-05 09:00:00',
  },
];
```

---

## 21.2 mockTags

```ts
export const mockTags: ContentTag[] = [
  {
    tagId: 1,
    tagName: 'Go语言',
    description: 'Go 后端开发、微服务、工程化实践。',
    status: 'enabled',
    useCount: 128,
    createdBy: 10003,
    createdAt: '2026-07-01 10:00:00',
    updatedAt: '2026-07-01 10:00:00',
  },
  {
    tagId: 2,
    tagName: 'React',
    description: 'React 前端开发、组件设计、状态管理。',
    status: 'enabled',
    useCount: 96,
    createdBy: 10003,
    createdAt: '2026-07-01 10:00:00',
    updatedAt: '2026-07-01 10:00:00',
  },
  {
    tagId: 3,
    tagName: 'AI工具',
    description: 'AI 编程助手、智能体、效率工具。',
    status: 'enabled',
    useCount: 180,
    createdBy: 10003,
    createdAt: '2026-07-01 10:00:00',
    updatedAt: '2026-07-01 10:00:00',
  },
  {
    tagId: 4,
    tagName: '项目实战',
    description: '个人项目、MVP、工程实践。',
    status: 'enabled',
    useCount: 220,
    createdBy: 10003,
    createdAt: '2026-07-01 10:00:00',
    updatedAt: '2026-07-01 10:00:00',
  },
];
```

---

## 21.3 mockTopics

```ts
export const mockTopics: Topic[] = [
  {
    topicId: 1,
    name: 'AI编程助手实践',
    description: '分享 AI 编程助手在真实项目中的使用经验。',
    coverImage: 'https://images.unsplash.com/photo-1677442136019-21780ecad995',
    postCount: 560,
    participantCount: 2345,
    isOfficial: true,
    isRecommended: true,
    status: 'enabled',
    createdBy: 10003,
    createdAt: '2026-07-01 10:00:00',
    updatedAt: '2026-07-01 10:00:00',
  },
  {
    topicId: 2,
    name: '个人项目部署',
    description: '讨论 Docker、Nginx、CI/CD、云服务器部署经验。',
    coverImage: 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31',
    postCount: 320,
    participantCount: 1200,
    isOfficial: true,
    isRecommended: true,
    status: 'enabled',
    createdBy: 10003,
    createdAt: '2026-07-01 10:00:00',
    updatedAt: '2026-07-01 10:00:00',
  },
];
```

---

## 21.4 mockCircles

```ts
export const mockCircles: Circle[] = [
  {
    circleId: 1,
    name: 'Go 后端开发圈',
    avatar: 'https://api.dicebear.com/7.x/shapes/svg?seed=go-circle',
    description: '专注 Go、微服务、性能优化和工程实践。',
    category: '后端开发',
    tags: [
      { tagId: 1, tagName: 'Go语言' },
      { tagId: 4, tagName: '项目实战' },
    ],
    ownerId: 10003,
    owner: {
      userId: 10003,
      nickname: '社区管理员',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin',
    },
    joinType: 'direct',
    postPermission: 'all',
    memberCount: 12345,
    postCount: 3210,
    featuredPostCount: 120,
    isJoined: true,
    myRole: 'member',
    myStatus: 'normal',
    isRecommended: true,
    status: 'normal',
    rules: '禁止广告引流；禁止人身攻击；提问前请先搜索历史内容。',
    createdAt: '2026-07-01 10:00:00',
    updatedAt: '2026-07-05 10:00:00',
  },
  {
    circleId: 2,
    name: 'React 前端圈',
    avatar: 'https://api.dicebear.com/7.x/shapes/svg?seed=react-circle',
    description: '讨论 React、TypeScript、组件化和前端工程化。',
    category: '前端开发',
    tags: [
      { tagId: 2, tagName: 'React' },
      { tagId: 4, tagName: '项目实战' },
    ],
    ownerId: 10002,
    owner: {
      userId: 10002,
      nickname: '前端小李',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=frontend',
    },
    joinType: 'approval',
    postPermission: 'all',
    memberCount: 8800,
    postCount: 2100,
    featuredPostCount: 88,
    isJoined: false,
    isRecommended: true,
    status: 'normal',
    rules: '保持友好讨论，禁止无意义刷屏。',
    createdAt: '2026-07-01 10:00:00',
    updatedAt: '2026-07-05 10:00:00',
  },
];
```

---

## 21.5 mockPosts

```ts
export const mockPosts: Post[] = [
  {
    postId: 20001,
    postType: 'original',
    authorId: 10002,
    author: {
      userId: 10002,
      nickname: '前端小李',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=frontend',
      level: 4,
      levelName: '社区达人',
    },
    title: '如何高效学习一门新的编程语言？',
    content:
      '最近在学习 Go 语言，感觉语法并不难，难点在于工程实践。我的方法是先快速过一遍语法，然后直接做一个小项目，在项目中遇到问题再回头补知识点。',
    summary: '最近在学习 Go 语言，感觉语法并不难，难点在于工程实践...',
    images: [
      'https://images.unsplash.com/photo-1516321318423-f06f85e504b3',
      'https://images.unsplash.com/photo-1515879218367-8466d910aaa4',
    ],
    tags: [
      { tagId: 1, tagName: 'Go语言' },
      { tagId: 4, tagName: '项目实战' },
    ],
    topics: [
      { topicId: 2, name: '个人项目部署' },
    ],
    circleId: 1,
    circle: {
      circleId: 1,
      name: 'Go 后端开发圈',
      avatar: 'https://api.dicebear.com/7.x/shapes/svg?seed=go-circle',
      memberCount: 12345,
      postCount: 3210,
    },
    status: 'published',
    isTop: false,
    isFeatured: true,
    isSelected: false,
    viewCount: 1280,
    likeCount: 432,
    commentCount: 32,
    favoriteCount: 88,
    shareCount: 24,
    repostCount: 12,
    hotScore: 3200,
    liked: false,
    favorited: false,
    followedAuthor: false,
    createdAt: '2026-07-05 09:30:00',
    publishedAt: '2026-07-05 09:35:00',
    updatedAt: '2026-07-05 09:35:00',
  },
  {
    postId: 20002,
    postType: 'repost',
    authorId: 10001,
    author: {
      userId: 10001,
      nickname: '张三',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan',
      level: 3,
      levelName: '内容贡献者',
    },
    title: '转发：如何高效学习一门新的编程语言？',
    content: '',
    summary: '这个学习方法很适合新手，先做项目再补知识点。',
    images: [],
    tags: [
      { tagId: 1, tagName: 'Go语言' },
    ],
    topics: [
      { topicId: 2, name: '个人项目部署' },
    ],
    status: 'published',
    isTop: false,
    isFeatured: false,
    isSelected: false,
    sourcePostId: 20001,
    repostComment: '这个学习方法很适合新手，先做项目再补知识点。',
    viewCount: 300,
    likeCount: 80,
    commentCount: 12,
    favoriteCount: 10,
    shareCount: 5,
    repostCount: 2,
    hotScore: 600,
    liked: false,
    favorited: false,
    followedAuthor: true,
    createdAt: '2026-07-05 11:00:00',
    publishedAt: '2026-07-05 11:00:00',
    updatedAt: '2026-07-05 11:00:00',
  },
];
```

---

## 21.6 mockComments

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
      level: 3,
      levelName: '内容贡献者',
    },
    content: '这个学习方法挺实用的，先做项目再补知识点确实效率高。',
    likeCount: 12,
    liked: false,
    status: 'normal',
    replies: [],
    createdAt: '2026-07-05 10:10:00',
    updatedAt: '2026-07-05 10:10:00',
  },
];
```

---

## 21.7 mockNotifications

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
    category: 'interaction',
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
    type: 'circle_join_approved',
    category: 'circle',
    title: '圈子加入申请已通过',
    content: '你的申请已通过，欢迎加入 Go 后端开发圈。',
    targetType: 'circle',
    targetId: 1,
    readStatus: 'unread',
    createdAt: '2026-07-05 11:00:00',
  },
];
```

---

## 21.8 mockGrowth

```ts
export const mockTasks: Task[] = [
  {
    taskId: 1,
    title: '完善个人资料',
    description: '上传头像并填写个人简介。',
    type: 'newbie',
    rewardPoints: 20,
    targetValue: 1,
    currentValue: 0,
    status: 'todo',
    actionText: '去完成',
    actionUrl: '/users/me',
  },
  {
    taskId: 2,
    title: '每日签到',
    description: '每天签到获得积分。',
    type: 'daily',
    rewardPoints: 5,
    targetValue: 1,
    currentValue: 1,
    status: 'done',
    actionText: '领取奖励',
    actionUrl: '/tasks',
  },
];

export const mockBadges: Badge[] = [
  {
    badgeId: 1,
    name: '首次发帖',
    icon: '🥇',
    description: '发布第一篇帖子后获得。',
    condition: '发布第一篇帖子',
    status: 'enabled',
    createdAt: '2026-07-01 10:00:00',
  },
  {
    badgeId: 2,
    name: '人气作者',
    icon: '🔥',
    description: '累计获得 1000 个点赞。',
    condition: '累计获得 1000 个点赞',
    status: 'enabled',
    createdAt: '2026-07-01 10:00:00',
  },
];
```

---

## 21.9 mockActivities

```ts
export const mockActivities: Activity[] = [
  {
    activityId: 1,
    type: 'topic',
    title: '本周讨论主题：AI 编程助手实践',
    description: '分享 AI 编程助手在真实项目中的使用经验。',
    coverImage: 'https://images.unsplash.com/photo-1677442136019-21780ecad995',
    topicId: 1,
    topicName: 'AI编程助手实践',
    startAt: '2026-07-01 00:00:00',
    endAt: '2026-07-15 23:59:59',
    participantCount: 2345,
    submissionCount: 560,
    status: 'ongoing',
  },
  {
    activityId: 2,
    type: 'vote',
    title: '你最常用的前端框架是？',
    description: '选择你当前最常用的前端框架。',
    coverImage: 'https://images.unsplash.com/photo-1555066931-4365d14bab8c',
    startAt: '2026-07-01 00:00:00',
    endAt: '2026-07-10 23:59:59',
    participantCount: 1024,
    submissionCount: 1024,
    status: 'ongoing',
    voteConfig: {
      question: '你最常用的前端框架是？',
      multiple: false,
      voted: false,
      selectedOptionIds: [],
      options: [
        { optionId: 1, text: 'React', voteCount: 560, percent: 55 },
        { optionId: 2, text: 'Vue', voteCount: 388, percent: 38 },
        { optionId: 3, text: 'Angular', voteCount: 50, percent: 5 },
        { optionId: 4, text: 'Svelte', voteCount: 26, percent: 2 },
      ],
    },
  },
];
```

---

# 22. API 通用约定

## 22.1 基础路径

```text
/api/v1
```

## 22.2 通用响应

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 22.3 分页响应

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

## 22.4 错误响应

```json
{
  "code": 40001,
  "message": "账号或密码错误",
  "data": null
}
```

## 22.5 错误码

|   错误码 | 说明      |
| ----: | ------- |
|     0 | 成功      |
| 40000 | 请求参数错误  |
| 40001 | 账号或密码错误 |
| 40100 | 未登录     |
| 40300 | 无权限     |
| 40301 | 用户已被禁言  |
| 40302 | 用户已被封禁  |
| 40303 | 圈子禁言    |
| 40304 | 圈子无发帖权限 |
| 40400 | 资源不存在   |
| 40900 | 数据冲突    |
| 50000 | 系统异常    |

## 22.6 请求头

```text
Authorization: Bearer <token>
```

---

# 23. 认证 API

## 23.1 登录

```text
POST /api/v1/auth/login
```

请求体：

```json
{
  "account": "zhangsan",
  "password": "123456"
}
```

响应：

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

## 23.2 注册

```text
POST /api/v1/auth/register
```

请求体：

```json
{
  "account": "lisi",
  "nickname": "李四",
  "password": "123456",
  "confirmPassword": "123456"
}
```

## 23.3 当前用户

```text
GET /api/v1/auth/me
```

## 23.4 退出登录

```text
POST /api/v1/auth/logout
```

---

# 24. 用户 API

```text
GET /api/v1/users/{userId}
PUT /api/v1/users/me/profile
GET /api/v1/users/me/posts
GET /api/v1/users/me/reposts
GET /api/v1/users/me/liked-posts
GET /api/v1/users/me/favorite-posts
GET /api/v1/users/me/following
GET /api/v1/users/me/circles
GET /api/v1/users/me/badges
GET /api/v1/users/{userId}/posts
```

## 24.1 更新资料

```text
PUT /api/v1/users/me/profile
```

请求体：

```json
{
  "nickname": "新的昵称",
  "avatar": "https://xxx.com/avatar.png",
  "bio": "新的个人简介"
}
```

---

# 25. 帖子 API

## 25.1 获取内容流

```text
GET /api/v1/posts
```

Query：

| 参数        | 类型     | 说明                                                   |
| --------- | ------ | ---------------------------------------------------- |
| page      | number | 页码                                                   |
| pageSize  | number | 每页数量                                                 |
| feedType  | string | recommend / latest / hot / following                 |
| sort      | string | recommend / latest / hot / comment / favorite / view |
| timeRange | string | all / today / week / month                           |
| tagId     | number | 标签 ID                                                |
| circleId  | number | 圈子 ID                                                |
| topicId   | number | 话题 ID                                                |
| keyword   | string | 搜索关键词                                                |

## 25.2 获取帖子详情

```text
GET /api/v1/posts/{postId}
```

## 25.3 创建帖子

```text
POST /api/v1/posts
```

请求体：

```json
{
  "title": "如何高效学习一门新的编程语言？",
  "content": "正文内容",
  "images": ["https://xxx.com/1.png"],
  "tagIds": [1, 4],
  "topicIds": [2],
  "circleId": 1,
  "publishMode": "now",
  "scheduledAt": null
}
```

## 25.4 更新帖子

```text
PUT /api/v1/posts/{postId}
```

## 25.5 删除帖子

```text
DELETE /api/v1/posts/{postId}
```

## 25.6 取消定时发布

```text
PUT /api/v1/posts/{postId}/cancel-schedule
```

## 25.7 立即发布定时帖子

```text
PUT /api/v1/posts/{postId}/publish-now
```

## 25.8 转发帖子

```text
POST /api/v1/posts/{postId}/repost
```

请求体：

```json
{
  "repostComment": "这个学习方法很适合新手。"
}
```

## 25.9 相似内容

```text
GET /api/v1/posts/{postId}/similar
```

## 25.10 圈子帖子置顶 / 取消置顶

```text
PUT /api/v1/posts/{postId}/top
DELETE /api/v1/posts/{postId}/top
```

## 25.11 圈子帖子加精 / 取消加精

```text
PUT /api/v1/posts/{postId}/feature
DELETE /api/v1/posts/{postId}/feature
```

---

# 26. 评论 API

```text
GET /api/v1/posts/{postId}/comments
POST /api/v1/comments
POST /api/v1/comments/{commentId}/replies
DELETE /api/v1/comments/{commentId}
```

## 26.1 发表评论

```json
{
  "postId": 20001,
  "content": "这个学习方法挺实用的。"
}
```

## 26.2 回复评论

```json
{
  "postId": 20001,
  "rootId": 30001,
  "replyToUserId": 10001,
  "content": "确实如此。"
}
```

---

# 27. 互动 API

```text
POST /api/v1/posts/{postId}/like
DELETE /api/v1/posts/{postId}/like

POST /api/v1/posts/{postId}/favorite
DELETE /api/v1/posts/{postId}/favorite

POST /api/v1/comments/{commentId}/like
DELETE /api/v1/comments/{commentId}/like

POST /api/v1/posts/{postId}/share
```

分享接口请求体：

```json
{
  "channel": "copy_link"
}
```

---

# 28. 关注 API

```text
POST /api/v1/users/{userId}/follow
DELETE /api/v1/users/{userId}/follow
POST /api/v1/users/{userId}/block
```

---

# 29. 标签 API

```text
GET /api/v1/tags
GET /api/v1/tags/{tagId}
GET /api/v1/tags/{tagId}/posts
```

## 29.1 获取标签列表

```text
GET /api/v1/tags?keyword=Go&status=enabled&page=1&pageSize=20
```

## 29.2 标签聚合内容

```text
GET /api/v1/tags/{tagId}/posts
```

Query：

| 参数        | 类型     | 说明                         |
| --------- | ------ | -------------------------- |
| page      | number | 页码                         |
| pageSize  | number | 每页数量                       |
| sort      | string | latest / hot               |
| timeRange | string | all / today / week / month |
| circleId  | number | 圈子 ID                      |
| topicId   | number | 话题 ID                      |

---

# 30. 圈子 API

## 30.1 圈子列表

```text
GET /api/v1/circles
```

Query：

| 参数       | 类型     | 说明                                 |
| -------- | ------ | ---------------------------------- |
| keyword  | string | 搜索关键词                              |
| category | string | 分类                                 |
| tagId    | number | 标签                                 |
| sort     | string | recommend / latest / member / post |
| page     | number | 页码                                 |
| pageSize | number | 每页数量                               |

## 30.2 创建圈子

```text
POST /api/v1/circles
```

请求体：

```json
{
  "name": "Go 后端开发圈",
  "avatar": "https://xxx.com/avatar.png",
  "description": "专注 Go、微服务、性能优化和工程实践。",
  "category": "后端开发",
  "tagIds": [1, 4],
  "joinType": "direct",
  "postPermission": "all",
  "rules": "禁止广告引流；禁止人身攻击。"
}
```

## 30.3 圈子详情

```text
GET /api/v1/circles/{circleId}
```

## 30.4 加入圈子

```text
POST /api/v1/circles/{circleId}/join
```

直接加入请求体：

```json
{}
```

审批加入请求体：

```json
{
  "reason": "想学习 Go 后端工程实践。"
}
```

## 30.5 退出圈子

```text
DELETE /api/v1/circles/{circleId}/join
```

## 30.6 圈子帖子

```text
GET /api/v1/circles/{circleId}/posts
```

Query：

| 参数       | 类型     | 说明                            |
| -------- | ------ | ----------------------------- |
| sort     | string | latest / hot / featured / top |
| tagId    | number | 标签                            |
| page     | number | 页码                            |
| pageSize | number | 每页数量                          |

## 30.7 圈子公告

```text
GET /api/v1/circles/{circleId}/announcements
POST /api/v1/circles/{circleId}/announcements
PUT /api/v1/circles/{circleId}/announcements/{announcementId}
DELETE /api/v1/circles/{circleId}/announcements/{announcementId}
```

## 30.8 圈子规则

```text
GET /api/v1/circles/{circleId}/rules
PUT /api/v1/circles/{circleId}/rules
```

## 30.9 圈子成员

```text
GET /api/v1/circles/{circleId}/members
```

## 30.10 圈子成员管理

```text
PUT /api/v1/circles/{circleId}/members/{userId}/role
PUT /api/v1/circles/{circleId}/members/{userId}/mute
PUT /api/v1/circles/{circleId}/members/{userId}/unmute
DELETE /api/v1/circles/{circleId}/members/{userId}
```

修改角色请求体：

```json
{
  "role": "moderator"
}
```

禁言请求体：

```json
{
  "durationDays": 7,
  "reason": "刷屏"
}
```

## 30.11 加入申请管理

```text
GET /api/v1/circles/{circleId}/join-requests
PUT /api/v1/circles/{circleId}/join-requests/{requestId}/approve
PUT /api/v1/circles/{circleId}/join-requests/{requestId}/reject
```

拒绝请求体：

```json
{
  "reason": "申请理由不充分"
}
```

---

# 31. 话题 API

```text
GET /api/v1/topics
GET /api/v1/topics/{topicId}
GET /api/v1/topics/{topicId}/posts
```

## 31.1 话题列表

```text
GET /api/v1/topics?type=hot&page=1&pageSize=20
```

## 31.2 话题内容

```text
GET /api/v1/topics/{topicId}/posts
```

Query：

| 参数        | 类型     | 说明                         |
| --------- | ------ | -------------------------- |
| sort      | string | latest / hot / featured    |
| timeRange | string | all / today / week / month |
| tagId     | number | 标签                         |
| page      | number | 页码                         |
| pageSize  | number | 每页数量                       |

---

# 32. 搜索 API

## 32.1 搜索联想

```text
GET /api/v1/search/suggest?keyword=Go
```

响应：

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "type": "post",
      "id": 20001,
      "title": "Go 后端工程化实践",
      "description": "帖子"
    }
  ]
}
```

## 32.2 搜索结果

```text
GET /api/v1/search
```

Query：

| 参数       | 类型     | 说明                                           |
| -------- | ------ | -------------------------------------------- |
| keyword  | string | 搜索词                                          |
| type     | string | all / post / comment / user / topic / circle |
| sort     | string | recommend / latest / hot                     |
| page     | number | 页码                                           |
| pageSize | number | 每页数量                                         |

## 32.3 热门搜索

```text
GET /api/v1/search/hot-keywords
```

---

# 33. 推荐 API

```text
GET /api/v1/recommend/posts
GET /api/v1/recommend/circles
GET /api/v1/recommend/topics
GET /api/v1/recommend/cold-start
```

## 33.1 推荐帖子

```text
GET /api/v1/recommend/posts?page=1&pageSize=10
```

## 33.2 冷启动推荐

```text
GET /api/v1/recommend/cold-start
```

返回热门帖子、推荐圈子、推荐话题。

---

# 34. 通知 API

```text
GET /api/v1/notifications
PUT /api/v1/notifications/{notificationId}/read
PUT /api/v1/notifications/read-all
GET /api/v1/notifications/unread-count
```

Query：

| 参数       | 类型     | 说明                                                                        |
| -------- | ------ | ------------------------------------------------------------------------- |
| category | string | all / interaction / follow / circle / review / growth / activity / system |
| page     | number | 页码                                                                        |
| pageSize | number | 每页数量                                                                      |

---

# 35. 举报 API

```text
POST /api/v1/reports
```

请求体：

```json
{
  "targetType": "post",
  "targetId": 20001,
  "reason": "spam",
  "description": "疑似广告内容"
}
```

---

# 36. 成长体系 API

## 36.1 签到

```text
POST /api/v1/growth/check-in
```

响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "checkedIn": true,
    "points": 5,
    "continuousDays": 3
  }
}
```

## 36.2 任务列表

```text
GET /api/v1/growth/tasks
```

Query：

| 参数   | 类型     | 说明                      |
| ---- | ------ | ----------------------- |
| type | string | newbie / daily / growth |

## 36.3 领取任务奖励

```text
POST /api/v1/growth/tasks/{taskId}/claim
```

## 36.4 用户勋章

```text
GET /api/v1/growth/users/{userId}/badges
```

## 36.5 排行榜

```text
GET /api/v1/growth/rankings
```

Query：

| 参数    | 类型     | 说明                                       |
| ----- | ------ | ---------------------------------------- |
| type  | string | active / contribution / creator / circle |
| range | string | day / week / all                         |

## 36.6 积分流水

```text
GET /api/v1/growth/point-logs
```

---

# 37. 活动 API

## 37.1 活动列表

```text
GET /api/v1/activities
```

Query：

| 参数     | 类型     | 说明                                   |
| ------ | ------ | ------------------------------------ |
| type   | string | all / topic / vote / essay / checkin |
| status | string | ongoing / ended / not_started        |

## 37.2 活动详情

```text
GET /api/v1/activities/{activityId}
```

## 37.3 参与活动

```text
POST /api/v1/activities/{activityId}/join
```

## 37.4 投票

```text
POST /api/v1/activities/{activityId}/vote
```

请求体：

```json
{
  "optionIds": [1]
}
```

## 37.5 打卡

```text
POST /api/v1/activities/{activityId}/check-in
```

请求体：

```json
{
  "content": "今天学习了 Go 的 context 使用。"
}
```

## 37.6 活动投稿列表

```text
GET /api/v1/activities/{activityId}/submissions
```

---

# 38. 官方公告 API

```text
GET /api/v1/announcements
GET /api/v1/announcements/{announcementId}
```

---

# 39. 上传 API

```text
POST /api/v1/upload/image
```

Content-Type：

```text
multipart/form-data
```

响应：

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

# 40. 后台 API

## 40.1 后台数据看板

```text
GET /api/v1/admin/dashboard/stats
GET /api/v1/admin/dashboard/user-trend
GET /api/v1/admin/dashboard/post-trend
GET /api/v1/admin/dashboard/activity-trend
GET /api/v1/admin/operation-dashboard
```

---

## 40.2 后台用户管理

```text
GET /api/v1/admin/users
PUT /api/v1/admin/users/{userId}/mute
PUT /api/v1/admin/users/{userId}/ban
PUT /api/v1/admin/users/{userId}/unban
```

---

## 40.3 后台帖子管理

```text
GET /api/v1/admin/posts
PUT /api/v1/admin/posts/{postId}/approve
PUT /api/v1/admin/posts/{postId}/reject
PUT /api/v1/admin/posts/{postId}/takedown
DELETE /api/v1/admin/posts/{postId}
PUT /api/v1/admin/posts/{postId}/select
DELETE /api/v1/admin/posts/{postId}/select
```

拒绝请求体：

```json
{
  "reason": "内容不符合社区规范"
}
```

下架请求体：

```json
{
  "reason": "内容存在违规信息"
}
```

---

## 40.4 后台评论管理

```text
GET /api/v1/admin/comments
DELETE /api/v1/admin/comments/{commentId}
```

---

## 40.5 后台审核管理

```text
GET /api/v1/admin/reviews/posts
GET /api/v1/admin/reviews/comments
GET /api/v1/admin/reviews/circles
PUT /api/v1/admin/reviews/{targetType}/{targetId}/approve
PUT /api/v1/admin/reviews/{targetType}/{targetId}/reject
```

---

## 40.6 后台举报管理

```text
GET /api/v1/admin/reports
PUT /api/v1/admin/reports/{reportId}/approve
PUT /api/v1/admin/reports/{reportId}/reject
```

通过举报请求体：

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
close_circle
```

---

## 40.7 后台标签管理

```text
GET /api/v1/admin/tags
POST /api/v1/admin/tags
PUT /api/v1/admin/tags/{tagId}
DELETE /api/v1/admin/tags/{tagId}
PUT /api/v1/admin/tags/{tagId}/enable
PUT /api/v1/admin/tags/{tagId}/disable
```

---

## 40.8 后台圈子管理

```text
GET /api/v1/admin/circles
GET /api/v1/admin/circles/{circleId}
PUT /api/v1/admin/circles/{circleId}/recommend
DELETE /api/v1/admin/circles/{circleId}/recommend
PUT /api/v1/admin/circles/{circleId}/close
PUT /api/v1/admin/circles/{circleId}/restore
DELETE /api/v1/admin/circles/{circleId}
```

关闭圈子请求体：

```json
{
  "reason": "圈子存在违规内容"
}
```

---

## 40.9 后台话题管理

```text
GET /api/v1/admin/topics
POST /api/v1/admin/topics
PUT /api/v1/admin/topics/{topicId}
DELETE /api/v1/admin/topics/{topicId}
PUT /api/v1/admin/topics/{topicId}/recommend
DELETE /api/v1/admin/topics/{topicId}/recommend
PUT /api/v1/admin/topics/{topicId}/enable
PUT /api/v1/admin/topics/{topicId}/disable
```

---

## 40.10 后台搜索管理

```text
GET /api/v1/admin/search/hot-keywords
PUT /api/v1/admin/search/hot-keywords/{keywordId}/block
PUT /api/v1/admin/search/hot-keywords/{keywordId}/restore
GET /api/v1/admin/search/trend
```

---

## 40.11 后台推荐管理

```text
GET /api/v1/admin/recommend-pool
POST /api/v1/admin/recommend-pool
PUT /api/v1/admin/recommend-pool/{id}
DELETE /api/v1/admin/recommend-pool/{id}
PUT /api/v1/admin/recommend-pool/{id}/online
PUT /api/v1/admin/recommend-pool/{id}/offline
```

---

## 40.12 后台积分等级管理

```text
GET /api/v1/admin/points/rules
PUT /api/v1/admin/points/rules/{ruleId}

GET /api/v1/admin/levels
POST /api/v1/admin/levels
PUT /api/v1/admin/levels/{level}

GET /api/v1/admin/points/logs
```

---

## 40.13 后台勋章管理

```text
GET /api/v1/admin/badges
POST /api/v1/admin/badges
PUT /api/v1/admin/badges/{badgeId}
DELETE /api/v1/admin/badges/{badgeId}
POST /api/v1/admin/badges/{badgeId}/grant
```

手动发放勋章请求体：

```json
{
  "userId": 10001
}
```

---

## 40.14 后台任务管理

```text
GET /api/v1/admin/tasks
POST /api/v1/admin/tasks
PUT /api/v1/admin/tasks/{taskId}
DELETE /api/v1/admin/tasks/{taskId}
PUT /api/v1/admin/tasks/{taskId}/enable
PUT /api/v1/admin/tasks/{taskId}/disable
GET /api/v1/admin/tasks/{taskId}/stats
```

---

## 40.15 后台活动管理

```text
GET /api/v1/admin/activities
POST /api/v1/admin/activities
PUT /api/v1/admin/activities/{activityId}
DELETE /api/v1/admin/activities/{activityId}
PUT /api/v1/admin/activities/{activityId}/online
PUT /api/v1/admin/activities/{activityId}/offline
GET /api/v1/admin/activities/{activityId}/stats
```

---

## 40.16 后台推荐位管理

```text
GET /api/v1/admin/banners
POST /api/v1/admin/banners
PUT /api/v1/admin/banners/{bannerId}
DELETE /api/v1/admin/banners/{bannerId}
PUT /api/v1/admin/banners/{bannerId}/enable
PUT /api/v1/admin/banners/{bannerId}/disable
```

---

## 40.17 后台官方公告管理

```text
GET /api/v1/admin/announcements
POST /api/v1/admin/announcements
PUT /api/v1/admin/announcements/{announcementId}
DELETE /api/v1/admin/announcements/{announcementId}
PUT /api/v1/admin/announcements/{announcementId}/publish
PUT /api/v1/admin/announcements/{announcementId}/offline
```

---

## 40.18 后台敏感词管理

```text
GET /api/v1/admin/sensitive-words
POST /api/v1/admin/sensitive-words
PUT /api/v1/admin/sensitive-words/{wordId}
DELETE /api/v1/admin/sensitive-words/{wordId}
PUT /api/v1/admin/sensitive-words/{wordId}/enable
PUT /api/v1/admin/sensitive-words/{wordId}/disable
```

---

## 40.19 后台操作日志

```text
GET /api/v1/admin/operation-logs
```

Query：

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

# 41. 前端状态管理建议

## 41.1 authStore

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

## 41.2 appStore

```ts
interface AppState {
  unreadNotificationCount: number;
  globalLoading: boolean;
  setUnreadNotificationCount: (count: number) => void;
  setGlobalLoading: (loading: boolean) => void;
}
```

## 41.3 draftStore

```ts
interface DraftState {
  postDraft: Partial<CreatePostRequest> | null;
  setPostDraft: (draft: Partial<CreatePostRequest>) => void;
  clearPostDraft: () => void;
}
```

---

# 42. API 文件函数拆分建议

## 42.1 auth.ts

```ts
export function login(data: LoginRequest) {}
export function register(data: RegisterRequest) {}
export function getCurrentUser() {}
export function logout() {}
```

## 42.2 post.ts

```ts
export function getPostList(params: PostQuery) {}
export function getPostDetail(postId: number) {}
export function createPost(data: CreatePostRequest) {}
export function updatePost(postId: number, data: UpdatePostRequest) {}
export function deletePost(postId: number) {}
export function cancelSchedulePost(postId: number) {}
export function publishPostNow(postId: number) {}
export function repost(postId: number, data: RepostRequest) {}
export function getSimilarPosts(postId: number) {}
export function topPost(postId: number) {}
export function untopPost(postId: number) {}
export function featurePost(postId: number) {}
export function unfeaturePost(postId: number) {}
```

## 42.3 tag.ts

```ts
export function getTags(params?: TagQuery) {}
export function getTagDetail(tagId: number) {}
export function getTagPosts(tagId: number, params: PostQuery) {}
```

## 42.4 circle.ts

```ts
export function getCircles(params: CircleQuery) {}
export function createCircle(data: CreateCircleRequest) {}
export function getCircleDetail(circleId: number) {}
export function joinCircle(circleId: number, reason?: string) {}
export function leaveCircle(circleId: number) {}
export function getCirclePosts(circleId: number, params: any) {}
export function getCircleAnnouncements(circleId: number) {}
export function createCircleAnnouncement(circleId: number, data: any) {}
export function getCircleRules(circleId: number) {}
export function updateCircleRules(circleId: number, rules: string) {}
export function getCircleMembers(circleId: number, params: any) {}
export function updateCircleMemberRole(circleId: number, userId: number, role: string) {}
export function muteCircleMember(circleId: number, userId: number, data: any) {}
export function unmuteCircleMember(circleId: number, userId: number) {}
export function removeCircleMember(circleId: number, userId: number) {}
export function getCircleJoinRequests(circleId: number) {}
export function approveCircleJoinRequest(circleId: number, requestId: number) {}
export function rejectCircleJoinRequest(circleId: number, requestId: number, reason: string) {}
```

## 42.5 search.ts

```ts
export function getSearchSuggest(keyword: string) {}
export function search(params: any) {}
export function getHotSearchKeywords() {}
```

## 42.6 growth.ts

```ts
export function checkIn() {}
export function getTasks(type: TaskType) {}
export function claimTaskReward(taskId: number) {}
export function getUserBadges(userId: number) {}
export function getRankings(params: { type: RankingType; range: RankingTimeRange }) {}
export function getPointLogs(params: any) {}
```

## 42.7 activity.ts

```ts
export function getActivities(params: any) {}
export function getActivityDetail(activityId: number) {}
export function joinActivity(activityId: number) {}
export function voteActivity(activityId: number, optionIds: number[]) {}
export function checkInActivity(activityId: number, content: string) {}
export function getActivitySubmissions(activityId: number) {}
```

## 42.8 admin.ts

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

export function getAdminTags(params: any) {}
export function createAdminTag(data: any) {}
export function updateAdminTag(tagId: number, data: any) {}

export function getAdminCircles(params: any) {}
export function closeAdminCircle(circleId: number, reason: string) {}
export function restoreAdminCircle(circleId: number) {}

export function getAdminTopics(params: any) {}
export function createAdminTopic(data: any) {}
export function updateAdminTopic(topicId: number, data: any) {}

export function getAdminActivities(params: any) {}
export function createAdminActivity(data: any) {}
export function updateAdminActivity(activityId: number, data: any) {}

export function getOperationLogs(params: any) {}
```

---

# 43. 页面与接口对应关系

## 43.1 前台页面

| 页面      | 依赖接口                                                  |
| ------- | ----------------------------------------------------- |
| 首页增强版   | `/posts`、`/recommend/posts`、`/banners`、点赞、收藏、分享、转发    |
| 发布帖子页   | `/posts`、`/tags`、`/topics`、`/circles`、`/upload/image` |
| 帖子详情页   | `/posts/{postId}`、评论、相似内容、互动接口                        |
| 标签聚合页   | `/tags/{tagId}`、`/tags/{tagId}/posts`                 |
| 圈子列表页   | `/circles`、加入圈子                                       |
| 创建圈子页   | `/circles`、`/upload/image`                            |
| 圈子首页    | `/circles/{circleId}`、圈子帖子、公告、规则、成员                   |
| 圈子成员管理页 | 成员列表、角色、禁言、踢出、加入申请                                    |
| 话题列表页   | `/topics`                                             |
| 话题详情页   | `/topics/{topicId}`、话题帖子                              |
| 搜索结果页   | `/search`                                             |
| 通知中心    | `/notifications`                                      |
| 我的主页    | `/auth/me`、我的帖子、我的圈子、我的勋章                             |
| 任务中心    | `/growth/tasks`、领取奖励、签到                               |
| 排行榜     | `/growth/rankings`                                    |
| 创作者中心   | 用户创作数据接口                                              |
| 活动中心    | `/activities`                                         |
| 活动详情    | 活动详情、投票、打卡、投稿                                         |
| 官方公告    | `/announcements`                                      |

## 43.2 后台页面

| 页面     | 依赖接口                                  |
| ------ | ------------------------------------- |
| 标签管理   | `/admin/tags`                         |
| 圈子管理   | `/admin/circles`                      |
| 话题管理   | `/admin/topics`                       |
| 搜索管理   | `/admin/search/hot-keywords`          |
| 推荐管理   | `/admin/recommend-pool`               |
| 积分等级管理 | `/admin/points/rules`、`/admin/levels` |
| 勋章管理   | `/admin/badges`                       |
| 任务管理   | `/admin/tasks`                        |
| 活动管理   | `/admin/activities`                   |
| 推荐位管理  | `/admin/banners`                      |
| 官方公告管理 | `/admin/announcements`                |
| 运营数据看板 | `/admin/operation-dashboard`          |

---

# 44. Mock 工具函数建议

## 44.1 模拟延迟

```ts
export function mockDelay<T>(data: T, delay = 300): Promise<T> {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(data);
    }, delay);
  });
}
```

## 44.2 模拟分页

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

## 44.3 模拟点赞

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

## 44.4 模拟收藏

```ts
export function togglePostFavorite(postId: number) {
  const post = mockPosts.find((item) => item.postId === postId);

  if (!post) {
    throw new Error('帖子不存在');
  }

  post.favorited = !post.favorited;
  post.favoriteCount += post.favorited ? 1 : -1;

  return {
    favorited: post.favorited,
    favoriteCount: post.favoriteCount,
  };
}
```

## 44.5 模拟筛选

```ts
export function filterPosts(params: PostQuery): Post[] {
  let result = [...mockPosts];

  if (params.tagId) {
    result = result.filter((post) =>
      post.tags.some((tag) => tag.tagId === params.tagId),
    );
  }

  if (params.circleId) {
    result = result.filter((post) => post.circleId === params.circleId);
  }

  if (params.topicId) {
    result = result.filter((post) =>
      post.topics.some((topic) => topic.topicId === params.topicId),
    );
  }

  if (params.keyword) {
    result = result.filter((post) =>
      post.title.includes(params.keyword || ''),
    );
  }

  return result;
}
```

---

# 45. 前端实现优先级

## P0：核心增强能力

```text
1. 首页增强版
2. 发布帖子增强版
3. 帖子详情增强版
4. 标签聚合页
5. 圈子列表页
6. 圈子首页
7. 话题列表页
8. 话题详情页
9. 搜索结果页
10. 分享弹窗
11. 转发弹窗
```

## P1：社区组织能力

```text
1. 创建圈子页
2. 圈子成员管理页
3. 圈子公告
4. 圈子规则
5. 置顶内容
6. 精华内容
7. 通知中心增强版
8. 我的主页增强版
```

## P2：增长运营能力

```text
1. 任务中心
2. 签到
3. 排行榜
4. 创作者中心
5. 活动中心
6. 活动详情页
7. 官方公告页
```

## P3：后台增强能力

```text
1. 标签管理
2. 圈子管理
3. 话题管理
4. 搜索管理
5. 推荐管理
6. 积分等级管理
7. 勋章管理
8. 任务管理
9. 活动管理
10. 推荐位管理
11. 官方公告管理
12. 运营数据看板
```

---

# 46. 给智能编程助手的最终实现提示词

请根据以下要求实现增强版社区 MVP 前端项目：

```text
使用 React + TypeScript + React Router + Zustand + Axios + Ant Design + Vite 实现一个增强版社区 MVP 前端项目。

项目需要包含前台和后台两套布局。

前台页面包括：
1. 登录页
2. 注册页
3. 首页增强版
4. 热门榜单页
5. 发布帖子增强版
6. 编辑帖子页
7. 帖子详情增强版
8. 标签聚合页
9. 圈子列表页
10. 创建圈子页
11. 圈子首页
12. 圈子成员管理页
13. 话题列表页
14. 话题详情页
15. 搜索结果页
16. 通知中心增强版
17. 我的主页增强版
18. 任务中心
19. 排行榜
20. 创作者中心
21. 活动中心
22. 活动详情页
23. 官方公告页

后台页面包括：
1. 数据看板
2. 用户管理
3. 帖子管理
4. 评论管理
5. 审核管理
6. 举报管理
7. 标签管理
8. 圈子管理
9. 话题管理
10. 搜索管理
11. 推荐管理
12. 积分等级管理
13. 勋章管理
14. 任务管理
15. 活动管理
16. 推荐位管理
17. 官方公告管理
18. 敏感词管理
19. 操作日志
20. 运营数据看板

实现要求：
1. 按照本文档的目录结构组织代码。
2. 按照本文档的组件拆分实现页面。
3. 使用 Mock 数据模拟后端接口。
4. 所有列表支持 loading、empty、error、pagination。
5. 所有表单支持基础校验。
6. 点赞、收藏、关注、分享、转发、加入圈子、签到、领取奖励等操作需要即时更新 UI。
7. 后台表格操作需要二次确认弹窗，操作成功后刷新表格。
8. 页面风格统一：浅灰背景、白色卡片、蓝色主色、圆角、轻阴影。
9. 顶部搜索框需要支持热门搜索和搜索联想。
10. 发布帖子页必须支持标签、话题、圈子、定时发布。
11. 圈子首页必须支持内容、公告、精华、成员、规则 Tab。
12. 后台必须支持标签、圈子、话题、活动、推荐位等管理页面。
```
